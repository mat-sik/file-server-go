package limited

import (
	"errors"
	"io"
)

type Buffer struct {
	buffer []byte
	offset int
}

type Resettable interface {
	Reset()
}

func (b *Buffer) Reset() {
	b.buffer = b.buffer[:0]
	b.offset = 0
}

func (b *Buffer) empty() bool {
	return len(b.buffer) == b.offset
}

type ReadableLength interface {
	Len() int
}

// Len returns amount of ready to be read unread bytes.
func (b *Buffer) Len() int {
	return len(b.buffer) - b.offset
}

func (b *Buffer) cap() int {
	return cap(b.buffer)
}

func (b *Buffer) available() int {
	return cap(b.buffer) - len(b.buffer)
}

func (b *Buffer) Write(p []byte) (int, error) {
	toReserve := len(p)
	oldLen := b.staticGrow(toReserve)
	availableBuffer := b.buffer[oldLen:]
	return copy(availableBuffer, p), nil
}

// staticGrow tries to make space for n, if it can't, it tires to make as much space as possible.
func (b *Buffer) staticGrow(n int) int {
	oldLen := len(b.buffer)
	newLen := min(b.cap(), oldLen+n)
	b.buffer = b.buffer[:newLen]
	return oldLen
}

func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	if b.empty() {
		b.Reset()
		return 0, nil
	}
	toWrite := b.buffer[b.offset:]
	n, err := w.Write(toWrite)
	b.offset += n
	return int64(n), err
}

type SingleWriterTo interface {
	SingleWriteTo(io.Writer, int) (int, error)
}

// SingleWriteTo writes at most N bytes from the underlying buffer to the io.Writer in a single Write operation.
// Returns ErrNotEnoughBuffered, if the buffer has not enough data.
func (b *Buffer) SingleWriteTo(w io.Writer, n int) (int, error) {
	if b.Len() < n {
		return 0, ErrNotEnoughBuffered
	}
	toWriteBytes := b.buffer[b.offset : b.offset+n]
	n, err := w.Write(toWriteBytes)
	if err != nil {
		return 0, err
	}
	b.offset += n
	return n, err
}

var ErrNotEnoughBuffered = errors.New("buffer has not enough buffered data")

type SingleReaderFrom interface {
	SingleReadFrom(io.Reader) (int, error)
}

// SingleReadFrom read as much as possible in a single Reader.read() call.
func (b *Buffer) SingleReadFrom(r io.Reader) (int, error) {
	if b.empty() {
		b.Reset()
	}
	availableBuffer := b.buffer[b.offset:cap(b.buffer)]
	n, err := r.Read(availableBuffer)
	if err != nil {
		return 0, err
	}
	b.buffer = availableBuffer[:n]
	return n, nil
}

func (b *Buffer) Read(p []byte) (int, error) {
	if b.empty() {
		b.Reset()
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	toRead := b.buffer[b.offset:]
	read := copy(p, toRead)
	b.offset += read
	return read, nil
}

type ByteIterator interface {
	Next(n int) []byte
}

func (b *Buffer) Next(n int) []byte {
	if b.empty() {
		b.Reset()
		return nil
	}
	n = min(len(b.buffer), n)
	data := b.buffer[b.offset : b.offset+n]
	b.offset += n
	return data
}

type BufferedAtLeastNEnsurer interface {
	EnsureBufferedAtLeastN(reader io.Reader, n int) error
}

func (b *Buffer) EnsureBufferedAtLeastN(reader io.Reader, n int) error {
	if !b.hasSpace(n) {
		return ErrSmallBuffer
	}
	for b.Len() < n {
		if _, err := b.SingleReadFrom(reader); err != nil {
			return err
		}
	}
	return nil
}

var ErrSmallBuffer = errors.New("buffer is too small")

func (b *Buffer) hasSpace(n int) bool {
	if b.Len()+n > b.cap() {
		return false
	}
	if b.available() < n {
		b.compact()
	}
	return true
}

func (b *Buffer) compact() {
	unread := b.buffer[b.offset:]
	b.Reset()
	_, _ = b.Write(unread)
}

func NewBuffer(buffer []byte) *Buffer {
	return &Buffer{buffer: buffer}
}
