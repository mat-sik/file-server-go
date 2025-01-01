package transfer

import (
	"encoding/json"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/header"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
)

func (d MessageDispatcher) SendMessage(m message.Message) error {
	return sendMessage(d.Buffer, d.Conn, d.HeaderBuffer, m)
}

func (d MessageDispatcher) ReceiveMessage() (message.Message, error) {
	return receiveMessage(d.Buffer, d.Conn)
}

type Messenger interface {
	io.WriterTo
	io.Writer
	io.Reader
	limited.BufferedAtLeastNEnsurer
	limited.ByteIterator
	limited.Resettable
	limited.ReadableLength
}

func sendMessage(messenger Messenger, writer io.Writer, headerBuffer []byte, m message.Message) error {
	defer messenger.Reset()

	encoder := json.NewEncoder(messenger)
	if err := encoder.Encode(m); err != nil {
		return err
	}

	messageSize := uint32(messenger.Len())
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
	if _, err := messenger.WriteTo(writer); err != nil {
		return err
	}
	return nil
}

func receiveMessage(messenger Messenger, reader io.Reader) (message.Message, error) {
	if err := messenger.EnsureBufferedAtLeastN(reader, header.Size); err != nil {
		return nil, err
	}

	messageHeader := header.DecodeHeader(messenger)

	toRead := messageHeader.PayloadSize - uint32(messenger.Len())
	if err := messenger.EnsureBufferedAtLeastN(reader, int(toRead)); err != nil {
		return nil, err
	}

	m, err := message.TypeNameConverter(messageHeader.PayloadType)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(messenger)
	if err = decoder.Decode(m); err != nil {
		return nil, err
	}

	return m, nil
}
