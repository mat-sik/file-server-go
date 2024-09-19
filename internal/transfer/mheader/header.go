package mheader

import (
	"bytes"
	"encoding/binary"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
)

type MessageHeader struct {
	PayloadSize uint32
	PayloadType message.TypeName
}

func EncodeHeader(header MessageHeader, headerBuffer []byte) error {
	if err := encodeMessageSize(header.PayloadSize, headerBuffer); err != nil {
		return err
	}
	if err := encodeMessageType(header.PayloadType, headerBuffer[transfer.Uint32ByteSize:]); err != nil {
		return err
	}
	return nil
}

func DecodeHeader(buffer *bytes.Buffer) MessageHeader {
	payloadSize := binary.BigEndian.Uint32(buffer.Next(transfer.Uint32ByteSize))
	payloadType := message.TypeName(binary.BigEndian.Uint64(buffer.Next(transfer.Uint64ByteSize)))
	return MessageHeader{
		PayloadSize: payloadSize,
		PayloadType: payloadType,
	}
}

func encodeMessageSize(messageSize uint32, headerBuffer []byte) error {
	if cap(headerBuffer) < transfer.Uint32ByteSize {
		return transfer.ErrHeaderBufferTooSmall
	}
	binary.BigEndian.PutUint32(headerBuffer, messageSize)
	return nil
}

func encodeMessageType(messageType message.TypeName, headerBuffer []byte) error {
	if cap(headerBuffer) < transfer.Uint64ByteSize {
		return transfer.ErrHeaderBufferTooSmall
	}
	binary.BigEndian.PutUint64(headerBuffer, uint64(messageType))
	return nil
}
