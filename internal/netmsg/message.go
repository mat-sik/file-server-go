package netmsg

import (
	"encoding/json"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg/header"
	"github.com/mat-sik/file-server-go/internal/netmsg/limited"
	"io"
)

type messageBuffer interface {
	io.WriterTo
	io.Writer
	io.Reader
	limited.MinReader
	limited.ByteIterator
	limited.Resettable
	limited.ReadableLength
}

func sendMessage(mess message.Message, headerBuffer []byte, buffer messageBuffer, writer io.Writer) error {
	defer buffer.Reset()

	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(mess); err != nil {
		return err
	}

	messageSize := uint32(buffer.Len())
	messageType := mess.GetType()
	messageHeader := header.Header{
		PayloadSize: messageSize,
		PayloadType: messageType,
	}
	if err := header.EncodeHeader(messageHeader, headerBuffer); err != nil {
		return err
	}

	if _, err := writer.Write(headerBuffer); err != nil {
		return err
	}
	if _, err := buffer.WriteTo(writer); err != nil {
		return err
	}
	return nil
}

func receiveMessage(reader io.Reader, buffer messageBuffer) (message.Message, error) {
	if err := buffer.ReadMin(reader, header.Size); err != nil {
		return nil, err
	}

	messageHeader := header.DecodeHeader(buffer)

	toRead := messageHeader.PayloadSize - uint32(buffer.Len())
	if err := buffer.ReadMin(reader, int(toRead)); err != nil {
		return nil, err
	}

	mess, err := message.TypeNameConverter(messageHeader.PayloadType)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(buffer)
	if err = decoder.Decode(mess); err != nil {
		return nil, err
	}

	return mess, nil
}
