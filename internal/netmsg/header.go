package netmsg

import (
	"encoding/binary"
	"errors"
)

type header struct {
	payloadSize uint32
}

func encodeHeader(header header, buffer []byte) error {
	if err := validateHeaderBufferSize(buffer, uint32ByteSize); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(buffer, header.payloadSize)
	return nil
}

func decodeHeader(buffer []byte) (header, error) {
	if len(buffer) < headerSize {
		return header{}, errors.New("buffer has not enough bytes to decode header")
	}
	uint32ByteSlice := buffer[:uint32ByteSize]
	payloadSize := binary.BigEndian.Uint32(uint32ByteSlice)

	return header{
		payloadSize: payloadSize,
	}, nil
}

func validateHeaderBufferSize(headerBuffer []byte, size int) error {
	if cap(headerBuffer) < size {
		return errors.New("header buffer has not enough capacity")
	}
	return nil
}

const (
	uint32ByteSize = 4
	headerSize     = uint32ByteSize
)
