package transfer

import (
	"encoding/json"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"github.com/mat-sik/file-server-go/internal/transfer/messheader"
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
	header := messheader.MessageHeader{
		PayloadSize: messageSize,
		PayloadType: messageType,
	}
	if err := messheader.EncodeHeader(header, headerBuffer); err != nil {
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
	if err := buffer.EnsureBufferedAtLeastN(reader, messheader.HeaderSize); err != nil {
		return nil, err
	}

	header := messheader.DecodeHeader(buffer)

	toRead := header.PayloadSize - uint32(buffer.Len())
	if ok := buffer.PrepareSpace(int(toRead)); !ok {
		return nil, ErrTooBigMessage
	}

	if err := buffer.EnsureBufferedAtLeastN(reader, int(toRead)); err != nil {
		return nil, err
	}

	m, err := message.TypeNameConverter(header.PayloadType)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(buffer)
	if err = decoder.Decode(m); err != nil {
		return nil, err
	}

	return m, nil
}

var ErrTooBigMessage = errors.New("buffer is too small to fit the message")
