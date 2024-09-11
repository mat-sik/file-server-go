package transfer

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"github.com/mat-sik/file-server-go/internal/Message"
	"io"
)

func transfer(
	reader io.Reader,
	writer io.Writer,
	buffer *bytes.Buffer,
	toTransfer int,
) error {
	bufferCapacity := int64(buffer.Cap())
	written := 0
	for {
		if buffered := len(buffer.Bytes()); buffered > 0 {
			limit := min(buffered, toTransfer-written)
			n, err := writer.Write(buffer.Next(limit))
			if err != nil {
				return err
			}
			written += n
			if written == toTransfer {
				break
			}
			buffer.Reset()
		}
		limitedReader := io.LimitReader(reader, bufferCapacity)
		if _, err := buffer.ReadFrom(limitedReader); err != nil {
			return err
		}
	}
	return nil
}

func receiveMessage(
	reader io.Reader,
	buffer *bytes.Buffer,
) (*Message.Holder, error) {
	buffered := len(buffer.Bytes())
	if buffered < messageSizeByteAmount {
		if _, err := io.ReadAtLeast(reader, buffer.Bytes(), messageSizeByteAmount); err != nil {
			return nil, err
		}
		buffered = len(buffer.Bytes())
	}
	toRead := binary.BigEndian.Uint32(buffer.Next(messageSizeByteAmount)) - uint32(buffered)
	if err := ensureBufferHasSpace(buffer, toRead); err != nil {
		return nil, err
	}
	decoder := gob.NewDecoder(buffer)

	var message Message.Holder
	if err := decoder.Decode(&message); err != nil {
		return nil, err
	}

	return &message, nil
}

func ensureBufferHasSpace(buffer *bytes.Buffer, toRead uint32) error {
	bufferCapacity := uint32(buffer.Cap())
	buffered := uint32(len(buffer.Bytes()))
	if toRead+buffered > bufferCapacity {
		return ErrTooBigMessage
	}
	availableSize := uint32(buffer.Available())
	if availableSize < toRead {
		if err := compact(buffer); err != nil {
			return err
		}
	}
	return nil
}

func compact(buffer *bytes.Buffer) error {
	payload := buffer.Bytes()
	buffer.Reset()
	_, err := buffer.Write(payload)
	return err
}

var ErrTooBigMessage = errors.New("buffer is too small to fit the message")

const (
	messageSizeByteAmount = 4
)
