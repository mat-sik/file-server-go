package transfer

import (
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
