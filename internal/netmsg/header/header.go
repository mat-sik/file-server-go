package header

import (
	"encoding/binary"
	"errors"
)

type Header struct {
	PayloadSize uint32
}

func EncodeHeader(header Header, buffer []byte) error {
	if err := validateHeaderBufferSize(buffer, uint32ByteSize); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(buffer, header.PayloadSize)
	return nil
}

func DecodeHeader(buffer []byte) (Header, error) {
	if len(buffer) < Size {
		return Header{}, errors.New("buffer has not enough bytes to decode header")
	}
	uint32ByteSlice := buffer[:uint32ByteSize]
	payloadSize := binary.BigEndian.Uint32(uint32ByteSlice)

	return Header{
		PayloadSize: payloadSize,
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
	Size           = uint32ByteSize
)
