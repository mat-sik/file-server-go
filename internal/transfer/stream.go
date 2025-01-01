package transfer

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
)

func (d MessageDispatcher) StreamToNet(ctx context.Context, reader io.Reader, toTransfer int) error {
	return stream(ctx, reader, d.Conn, d.Buffer, toTransfer)
}

func (d MessageDispatcher) StreamFromNet(ctx context.Context, writer io.Writer, toTransfer int) error {
	return stream(ctx, d.Conn, writer, d.Buffer, toTransfer)
}

type Streamer interface {
	limited.SingleWriterTo
	limited.SingleReaderFrom
	limited.Resettable
	limited.ReadableLength
}

func stream(ctx context.Context, reader io.Reader, writer io.Writer, buffer Streamer, toTransfer int) error {
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
