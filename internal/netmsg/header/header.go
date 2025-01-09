package header

import (
	"encoding/binary"
	"errors"
	"github.com/mat-sik/file-server-go/internal/netmsg/limited"
)

type Header struct {
	PayloadSize uint32
}

func EncodeHeader(header Header, headerBuffer []byte) error {
	return encodeMessageSize(header.PayloadSize, headerBuffer)
}

func DecodeHeader(byteIterator limited.ByteIterator) Header {
	uint32ByteSlice := byteIterator.Next(uint32ByteSize)
	payloadSize := binary.BigEndian.Uint32(uint32ByteSlice)

	return Header{
		PayloadSize: payloadSize,
	}
}

func encodeMessageSize(messageSize uint32, headerBuffer []byte) error {
	if err := validateHeaderBufferSize(headerBuffer, uint32ByteSize); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(headerBuffer, messageSize)
	return nil
}

func validateHeaderBufferSize(headerBuffer []byte, size int) error {
	if cap(headerBuffer) < size {
		return errors.New("header buffer too small")
	}
	return nil
}

const (
	uint32ByteSize = 4
	Size           = uint32ByteSize
)
