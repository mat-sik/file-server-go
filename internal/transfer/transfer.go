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
