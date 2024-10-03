package transfer

import (
	"bytes"
	"context"
	"io"
)

func Stream(
	ctx context.Context,
	reader io.Reader,
	writer io.Writer,
	buffer *bytes.Buffer,
	toTransfer int,
) error {
	bufferCapacity := buffer.Cap()
	written := 0
	for {
		if err := ctxEarlyReturn(ctx); err != nil {
			return err
		}

		if buffered := buffer.Len(); buffered > 0 {
			toRead := toTransfer - written
			n, err := limitedWrite(buffer, writer, buffered, toRead)
			if err != nil {
				return err
			}
			written += n
			if written == toTransfer {
				break
			}
		}
		buffer.Reset()

		toRead := toTransfer - written
		limit := min(toRead, bufferCapacity)
		if _, err := limitedRead(buffer, reader, limit); err != nil {
			return err
		}
	}
	return nil
}

func limitedWrite(buffer *bytes.Buffer, writer io.Writer, buffered, toRead int) (int, error) {
	limit := min(buffered, toRead)
	toWriteBytes := buffer.Next(limit)
	return writer.Write(toWriteBytes)
}

func limitedRead(buffer *bytes.Buffer, reader io.Reader, limit int) (int64, error) {
	limitedReader := io.LimitReader(reader, int64(limit))
	return buffer.ReadFrom(limitedReader)
}

func ctxEarlyReturn(ctx context.Context) error {
	select {
	default:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
