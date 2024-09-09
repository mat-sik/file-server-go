package transfer

import (
	"bytes"
	"io"
)

func transfer(
	reader io.Reader,
	writer io.Writer,
	buffer *bytes.Buffer,
	toTransfer int,
) error {
	limitedReader := io.LimitedReader{R: reader, N: int64(buffer.Cap())}
	written := 0
	for {
		buffered := len(buffer.Bytes())
		if buffered > 0 {
			limit := min(buffered, toTransfer-written)
			n, err := writer.Write(buffer.Next(limit))
			if err != nil {
				panic(err)
			}
			written += n
			if written == toTransfer {
				break
			}
			buffer.Reset()
		}
		copiedLimitedReader := limitedReader
		if _, err := buffer.ReadFrom(&copiedLimitedReader); err != nil {
			panic(err)
		}
	}
	return nil
}
