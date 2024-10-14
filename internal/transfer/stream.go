package transfer

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
)

type StreamableBuffer interface {
	limited.SingleWriterTo
	limited.SingleReaderFrom
	Len() int
	Reset()
}

func Stream(ctx context.Context, reader io.Reader, writer io.Writer, buffer StreamableBuffer, toTransfer int) error {
	if toTransfer == 0 {
		return nil
	}

	written := 0
	for {
		if err := ctxEarlyReturn(ctx); err != nil {
			return err
		}

		if buffered := buffer.Len(); buffered > 0 {
			toRead := toTransfer - written
			limit := min(buffered, toRead)
			n, err := buffer.SingleWriteTo(writer, limit)
			if err != nil {
				return err
			}
			written += n
			if written == toTransfer {
				break
			}
		}
		buffer.Reset()

		if _, err := buffer.SingleReadFrom(reader); err != nil {
			return err
		}
	}
	return nil
}

func ctxEarlyReturn(ctx context.Context) error {
	select {
	default:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
