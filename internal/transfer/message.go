package transfer

import (
	"encoding/json"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/header"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
)

type Messenger interface {
	io.WriterTo
	io.Writer
	io.Reader
	limited.BufferedAtLeastNEnsurer
	limited.ByteIterator
	limited.Resettable
	limited.ReadableLength
}

func (dispatcher MessageDispatcher) SendMessage(
	m message.Message,
) error {
	defer dispatcher.Buffer.Reset()

	encoder := json.NewEncoder(dispatcher.Buffer)
	if err := encoder.Encode(m); err != nil {
		return err
	}

	messageSize := uint32(dispatcher.Buffer.Len())
	messageType := m.GetType()
	messageHeader := header.Header{
		PayloadSize: messageSize,
		PayloadType: messageType,
	}
	if err := header.EncodeHeader(messageHeader, dispatcher.HeaderBuffer); err != nil {
		return err
	}

	if _, err := dispatcher.ReadWriteCloser.Write(dispatcher.HeaderBuffer); err != nil {
		return err
	}
	if _, err := dispatcher.Buffer.WriteTo(dispatcher.ReadWriteCloser); err != nil {
		return err
	}
	return nil
}

func (dispatcher MessageDispatcher) ReceiveMessage() (message.Message, error) {
	if err := dispatcher.Buffer.EnsureBufferedAtLeastN(dispatcher.ReadWriteCloser, header.Size); err != nil {
		return nil, err
	}

	messageHeader := header.DecodeHeader(dispatcher.Buffer)

	toRead := messageHeader.PayloadSize - uint32(dispatcher.Buffer.Len())
	if err := dispatcher.Buffer.EnsureBufferedAtLeastN(dispatcher.ReadWriteCloser, int(toRead)); err != nil {
		return nil, err
	}

	m, err := message.TypeNameConverter(messageHeader.PayloadType)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(dispatcher.Buffer)
	if err = decoder.Decode(m); err != nil {
		return nil, err
	}

	return m, nil
}
