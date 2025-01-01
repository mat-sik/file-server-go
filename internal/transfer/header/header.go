package header

import (
	"encoding/binary"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
)

type Header struct {
	PayloadSize uint32
	PayloadType message.TypeName
}

func EncodeHeader(header Header, headerBuffer []byte) error {
	messageSize := header.PayloadSize
	if err := encodeMessageSize(messageSize, headerBuffer); err != nil {
		return err
	}
	messageType := header.PayloadType
	if err := encodeMessageType(messageType, headerBuffer[uint32ByteSize:]); err != nil {
		return err
	}
	return nil
}

func DecodeHeader(byteIterator limited.ByteIterator) Header {
	uint32ByteSlice := byteIterator.Next(uint32ByteSize)
	payloadSize := binary.BigEndian.Uint32(uint32ByteSlice)

	uint64ByteSlice := byteIterator.Next(uint64ByteSize)
	payloadType := message.TypeName(binary.BigEndian.Uint64(uint64ByteSlice))

	return Header{
		PayloadSize: payloadSize,
		PayloadType: payloadType,
	}
}

func encodeMessageSize(messageSize uint32, headerBuffer []byte) error {
	if err := validateHeaderBufferSize(headerBuffer, uint32ByteSize); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(headerBuffer, messageSize)
	return nil
}

func encodeMessageType(messageType message.TypeName, headerBuffer []byte) error {
	if err := validateHeaderBufferSize(headerBuffer, uint64ByteSize); err != nil {
		return err
	}
	binary.BigEndian.PutUint64(headerBuffer, uint64(messageType))
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
	uint64ByteSize = 8
	Size           = uint32ByteSize + uint64ByteSize
)
