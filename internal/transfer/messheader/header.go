package messheader

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
)

type MessageHeader struct {
	PayloadSize uint32
	PayloadType message.TypeName
}

func EncodeHeader(header MessageHeader, headerBuffer []byte) error {
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

func DecodeHeader(buffer *bytes.Buffer) MessageHeader {
	uint32ByteSlice := buffer.Next(uint32ByteSize)
	payloadSize := binary.BigEndian.Uint32(uint32ByteSlice)

	uint64ByteSlice := buffer.Next(uint64ByteSize)
	payloadType := message.TypeName(binary.BigEndian.Uint64(uint64ByteSlice))

	return MessageHeader{
		PayloadSize: payloadSize,
		PayloadType: payloadType,
	}
}

func encodeMessageSize(messageSize uint32, headerBuffer []byte) error {
	if cap(headerBuffer) < uint32ByteSize {
		return ErrHeaderBufferTooSmall
	}
	binary.BigEndian.PutUint32(headerBuffer, messageSize)
	return nil
}

func encodeMessageType(messageType message.TypeName, headerBuffer []byte) error {
	if cap(headerBuffer) < uint64ByteSize {
		return ErrHeaderBufferTooSmall
	}
	binary.BigEndian.PutUint64(headerBuffer, uint64(messageType))
	return nil
}

var ErrHeaderBufferTooSmall = errors.New("header buffer too small")

const (
	uint32ByteSize = 4
	uint64ByteSize = 8
	HeaderSize     = uint32ByteSize + uint64ByteSize
)
