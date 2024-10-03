package transfer

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/mheader"
	"io"
)

func SendMessage(
	writer io.Writer,
	headerBuffer []byte,
	messageBuffer *bytes.Buffer,
	m message.Message,
) error {
	defer messageBuffer.Reset()

	encoder := json.NewEncoder(messageBuffer)
	if err := encoder.Encode(m); err != nil {
		return err
	}

	messageSize := uint32(messageBuffer.Len())
	messageType := m.GetType()
	header := mheader.MessageHeader{
		PayloadSize: messageSize,
		PayloadType: messageType,
	}
	if err := mheader.EncodeHeader(header, headerBuffer); err != nil {
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
	buffer *bytes.Buffer,
) (message.Message, error) {
	if err := readN(reader, buffer, mheader.HeaderSize); err != nil {
		return nil, err
	}

	header := mheader.DecodeHeader(buffer)

	toRead := header.PayloadSize - uint32(buffer.Len())
	if err := ensureBufferHasSpace(buffer, toRead); err != nil {
		return nil, err
	}

	if err := readN(reader, buffer, int(toRead)); err != nil {
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

func readN(reader io.Reader, buffer *bytes.Buffer, n int) error {
	limit := int64(n)
	limitedReader := io.LimitReader(reader, limit)
	_, err := buffer.ReadFrom(limitedReader)
	return err
}

func ensureBufferHasSpace(buffer *bytes.Buffer, size uint32) error {
	bufferCapacity := uint32(buffer.Cap())
	buffered := uint32(buffer.Len())
	if size+buffered > bufferCapacity {
		return ErrTooBigMessage
	}
	availableSpace := uint32(buffer.Available())
	if availableSpace < size {
		if err := compact(buffer); err != nil {
			return err
		}
	}
	return nil
}

func compact(buffer *bytes.Buffer) error {
	payload := buffer.Bytes()
	buffer.Reset()
	_, err := buffer.Write(payload)
	return err
}

var ErrTooBigMessage = errors.New("buffer is too small to fit the message")
