package transfer

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
)

func Stream(
	ctx context.Context,
	reader io.Reader,
	writer io.Writer,
	buffer *limited.Buffer,
	toTransfer int,
) error {
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

		if _, err := buffer.MaxRead(reader); err != nil {
			return err
		}
	}
	return nil
}

func limitedWrite(buffer *limited.Buffer, writer io.Writer, buffered, toRead int) (int, error) {
	limit := min(buffered, toRead)
	toWriteBytes := buffer.Next(limit)
	return writer.Write(toWriteBytes)
}

func ctxEarlyReturn(ctx context.Context) error {
	select {
	default:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
