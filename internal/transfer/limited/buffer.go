package limited

import (
	"io"
)

type Buffer struct {
	buffer []byte
	offset int
}

func (b *Buffer) Reset() {
	b.buffer = b.buffer[:0]
	b.offset = 0
}

func (b *Buffer) empty() bool {
	return len(b.buffer) == b.offset
}

func (b *Buffer) Len() int {
	return len(b.buffer) - b.offset
}

func (b *Buffer) Cap() int {
	return cap(b.buffer)
}

func (b *Buffer) Available() int {
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
	newLen := min(b.Cap(), oldLen+n)
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

func (b *Buffer) LimitedWrite(w io.Writer, limit int) (int, error) {
	toWriteBytes := b.Next(limit)
	return w.Write(toWriteBytes)
}

// MaxRead read as much as possible in a single cycle
func (b *Buffer) MaxRead(r io.Reader) (int, error) {
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

func (b *Buffer) EnsureBufferedAtLeastN(reader io.Reader, n int) error {
	for b.Len() < n {
		if _, err := b.MaxRead(reader); err != nil {
			return err
		}
	}
	return nil
}

func (b *Buffer) PrepareSpace(n int) bool {
	if b.Len()+n > b.Cap() {
		return false
	}
	if b.Available() < n {
		b.Compact()
	}
	return true
}

func (b *Buffer) Compact() {
	unread := b.buffer[b.offset:]
	b.Reset()
	_, _ = b.Write(unread)
}

func NewBuffer(buffer []byte) *Buffer {
	return &Buffer{buffer: buffer}
}
