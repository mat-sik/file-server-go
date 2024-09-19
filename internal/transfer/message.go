package transfer

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/mat-sik/file-server-go/internal/message"
	"io"
)

func sendMessage(
	writer io.Writer,
	headerBuffer []byte,
	messageBuffer *bytes.Buffer,
	holder *message.Holder,
) error {
	encoder := json.NewEncoder(messageBuffer)
	if err := encoder.Encode(holder.PayloadStruct); err != nil {
		return err
	}

	messageSize := uint32(messageBuffer.Len())
	messageType := holder.PayloadType
	header := messageHeader{
		payloadSize: messageSize,
		payloadType: messageType,
	}
	if err := encodeHeader(header, headerBuffer); err != nil {
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

type messageHeader struct {
	payloadSize uint32
	payloadType message.TypeName
}

func encodeHeader(header messageHeader, headerBuffer []byte) error {
	if err := encodeMessageSize(header.payloadSize, headerBuffer); err != nil {
		return err
	}
	if err := encodeMessageType(header.payloadType, headerBuffer[Uint32ByteSize:]); err != nil {
		return err
	}
	return nil
}

func decodeHeader(buffer *bytes.Buffer) messageHeader {
	payloadSize := binary.BigEndian.Uint32(buffer.Next(Uint32ByteSize))
	payloadType := message.TypeName(binary.BigEndian.Uint64(buffer.Next(Uint64ByteSize)))
	return messageHeader{
		payloadSize: payloadSize,
		payloadType: payloadType,
	}
}

func encodeMessageSize(messageSize uint32, headerBuffer []byte) error {
	if cap(headerBuffer) < Uint32ByteSize {
		return ErrHeaderBufferTooSmall
	}
	binary.BigEndian.PutUint32(headerBuffer, messageSize)
	return nil
}

func encodeMessageType(messageType message.TypeName, headerBuffer []byte) error {
	if cap(headerBuffer) < Uint64ByteSize {
		return ErrHeaderBufferTooSmall
	}
	binary.BigEndian.PutUint64(headerBuffer, uint64(messageType))
	return nil
}

func receiveMessage(
	ctx context.Context,
	reader io.Reader,
	buffer *bytes.Buffer,
) (message.Holder, error) {
	if err := ensureBuffered(ctx, reader, buffer, HeaderSize); err != nil {
		return message.Holder{}, err
	}

	header := decodeHeader(buffer)

	toRead := header.payloadSize - uint32(buffer.Len())
	if err := ensureBufferHasSpace(buffer, toRead); err != nil {
		return message.Holder{}, err
	}

	payload, err := message.TypeNameConverter(header.payloadType)
	if err != nil {
		return message.Holder{}, err
	}

	decoder := json.NewDecoder(buffer)
	if err = decoder.Decode(payload); err != nil {
		return message.Holder{}, err
	}

	return message.Holder{
		PayloadType:   header.payloadType,
		PayloadStruct: payload,
	}, nil
}

func ensureBuffered(ctx context.Context, reader io.Reader, buffer *bytes.Buffer, min int) error {
	if buffer.Len() < Uint32ByteSize {
		if _, err := readAtLeast(ctx, reader, buffer, min); err != nil {
			return err
		}
	}
	return nil
}

func readAtLeast(ctx context.Context, reader io.Reader, buffer *bytes.Buffer, min int) (int, error) {
	for {
		select {
		default:
		case <-ctx.Done():
			return buffer.Len(), ctx.Err()
		}
		availableSpace := int64(buffer.Available())
		limitedReader := io.LimitReader(reader, availableSpace)
		if _, err := buffer.ReadFrom(limitedReader); err != nil {
			return buffer.Len(), err
		}
		if buffer.Len() >= min {
			return buffer.Len(), nil
		}
	}
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
