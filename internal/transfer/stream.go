package transfer

import (
	"bytes"
	"context"
	"io"
)

func stream(
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
		default:
		case <-ctx.Done():
			return ctx.Err()
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
