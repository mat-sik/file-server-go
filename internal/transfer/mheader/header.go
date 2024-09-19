package mheader

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
	if err := encodeMessageSize(header.PayloadSize, headerBuffer); err != nil {
		return err
	}
	if err := encodeMessageType(header.PayloadType, headerBuffer[uint32ByteSize:]); err != nil {
		return err
	}
	return nil
}

func DecodeHeader(buffer *bytes.Buffer) MessageHeader {
	payloadSize := binary.BigEndian.Uint32(buffer.Next(uint32ByteSize))
	payloadType := message.TypeName(binary.BigEndian.Uint64(buffer.Next(uint64ByteSize)))
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

var ErrHeaderBufferTooSmall = errors.New("mheader buffer too small")

const (
	uint32ByteSize = 4
	uint64ByteSize = 8
	HeaderSize     = uint32ByteSize + uint64ByteSize
)
