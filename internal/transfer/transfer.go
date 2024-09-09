package transfer

import (
	"bytes"
	"io"
)

func transfer(
	reader io.Reader,
	writer io.Writer,
	buffer []byte,
	offset int,
	buffered int,
	toTransfer int,
) (int, int, error) {
	written := 0
	for {
		if buffered > 0 {
			end := min(offset+buffered, toTransfer-written)
			n, err := writer.Write(buffer[offset:end])
			if err != nil {
				panic(err)
			}
			written += n
			buffered -= n
			if written == toTransfer {
				offset = n
				break
			}
		}
		offset = 0
		n, err := reader.Read(buffer)
		if err != nil {
			panic(err)
		}
		buffered = n
	}
	return offset, buffered, nil
}

func transferWithBuffer(
	reader io.Reader,
	writer io.Writer,
	buffer bytes.Buffer,
	toTransfer int,
) error {
	limitedReader := io.LimitedReader{R: reader, N: int64(buffer.Cap())}
	written := int64(0)
	for {
		if len(buffer.Bytes()) > 0 {
			n, err := buffer.WriteTo(writer)
			if err != nil {
				panic(err)
			}
			written += n
			if written == int64(toTransfer) {
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
