package transfer

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
	"io"
)

func transfer(
	ctx context.Context,
	reader io.Reader,
	writer io.Writer,
	buffer *bytes.Buffer,
	toTransfer int,
) error {
	bufferCapacity := int64(buffer.Cap())
	written := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if buffered := buffer.Len(); buffered > 0 {
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

func sendMessage(
	writer io.Writer,
	sizeBuffer []byte,
	messageBuffer *bytes.Buffer,
	holder *message.Holder,
) error {
	encoder := gob.NewEncoder(messageBuffer)
	if err := encoder.Encode(holder); err != nil {
		return err
	}
	encodeMessageSize(messageBuffer, sizeBuffer)
	if _, err := writer.Write(sizeBuffer); err != nil {
		return err
	}
	if _, err := messageBuffer.WriteTo(writer); err != nil {
		return err
	}
	return nil
}

func encodeMessageSize(messageBuffer *bytes.Buffer, sizeBuffer []byte) {
	encodedHolderSize := messageBuffer.Len()
	binary.BigEndian.PutUint32(sizeBuffer, uint32(encodedHolderSize))
}

func receiveMessage(
	ctx context.Context,
	reader io.Reader,
	buffer *bytes.Buffer,
) (message.Holder, error) {
	if err := ensureBuffered(ctx, reader, buffer, messageSizeByteAmount); err != nil {
		return message.Holder{}, err
	}
	toRead := binary.BigEndian.Uint32(buffer.Next(messageSizeByteAmount)) - uint32(buffer.Len())
	if err := ensureBufferHasSpace(buffer, toRead); err != nil {
		return message.Holder{}, err
	}
	decoder := gob.NewDecoder(buffer)

	var holder message.Holder
	if err := decoder.Decode(&holder); err != nil {
		return message.Holder{}, err
	}

	return holder, nil
}

func ensureBuffered(ctx context.Context, reader io.Reader, buffer *bytes.Buffer, min int) error {
	if buffer.Len() < messageSizeByteAmount {
		if _, err := readAtLeast(ctx, reader, buffer, min); err != nil {
			return err
		}
	}
	return nil
}

func readAtLeast(ctx context.Context, reader io.Reader, buffer *bytes.Buffer, min int) (int, error) {
	for {
		select {
		case <-ctx.Done():
			return buffer.Len(), ctx.Err()
		default:
		}
		availableSpace := int64(buffer.Available())
		limitedReader := io.LimitReader(reader, availableSpace)
		if _, err := buffer.ReadFrom(limitedReader); err != nil {
			return buffer.Len(), err
		}
		if buffer.Len() >= min {
			return buffer.Len(), nil
		}
	}
}

func ensureBufferHasSpace(buffer *bytes.Buffer, size uint32) error {
	bufferCapacity := uint32(buffer.Cap())
	buffered := uint32(buffer.Len())
	if size+buffered > bufferCapacity {
		return ErrTooBigMessage
	}
	availableSpace := uint32(buffer.Available())
	if availableSpace < size {
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
