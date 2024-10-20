package transfer

import (
	"encoding/json"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/header"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
)

func SendMessage(
	writer io.Writer,
	headerBuffer []byte,
	messageBuffer *limited.Buffer,
	m message.Message,
) error {
	defer messageBuffer.Reset()

	encoder := json.NewEncoder(messageBuffer)
	if err := encoder.Encode(m); err != nil {
		return err
	}

	messageSize := uint32(messageBuffer.Len())
	messageType := m.GetType()
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
	if _, err := messageBuffer.WriteTo(writer); err != nil {
		return err
	}
	return nil
}

func ReceiveMessage(
	reader io.Reader,
	buffer *limited.Buffer,
) (message.Message, error) {
	if err := buffer.EnsureBufferedAtLeastN(reader, header.Size); err != nil {
		return nil, err
	}

	messageHeader := header.DecodeHeader(buffer)

	toRead := messageHeader.PayloadSize - uint32(buffer.Len())
	if err := buffer.EnsureBufferedAtLeastN(reader, int(toRead)); err != nil {
		return nil, err
	}

	m, err := message.TypeNameConverter(messageHeader.PayloadType)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(buffer)
	if err = decoder.Decode(m); err != nil {
		return nil, err
	}

	return m, nil
}
