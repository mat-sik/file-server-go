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

func (d MessageDispatcher) SendMessage(m message.Message) error {
	defer d.Buffer.Reset()

	encoder := json.NewEncoder(d.Buffer)
	if err := encoder.Encode(m); err != nil {
		return err
	}

	messageSize := uint32(d.Buffer.Len())
	messageType := m.GetType()
	messageHeader := header.Header{
		PayloadSize: messageSize,
		PayloadType: messageType,
	}
	if err := header.EncodeHeader(messageHeader, d.HeaderBuffer); err != nil {
		return err
	}

	if _, err := d.Conn.Write(d.HeaderBuffer); err != nil {
		return err
	}
	if _, err := d.Buffer.WriteTo(d.Conn); err != nil {
		return err
	}
	return nil
}

func (d MessageDispatcher) ReceiveMessage() (message.Message, error) {
	if err := d.Buffer.EnsureBufferedAtLeastN(d.Conn, header.Size); err != nil {
		return nil, err
	}

	messageHeader := header.DecodeHeader(d.Buffer)

	toRead := messageHeader.PayloadSize - uint32(d.Buffer.Len())
	if err := d.Buffer.EnsureBufferedAtLeastN(d.Conn, int(toRead)); err != nil {
		return nil, err
	}

	m, err := message.TypeNameConverter(messageHeader.PayloadType)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(d.Buffer)
	if err = decoder.Decode(m); err != nil {
		return nil, err
	}

	return m, nil
}
