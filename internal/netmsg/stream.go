package netmsg

import (
	"context"
	"github.com/mat-sik/limbuf/limited"
	"io"
)

type streamBuffer interface {
	limited.SingleWriterTo
	limited.SingleReaderFrom
	limited.Resettable
	limited.ReadableLength
}

func stream(ctx context.Context, reader io.Reader, writer io.Writer, buffer streamBuffer, toTransfer int) error {
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
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
