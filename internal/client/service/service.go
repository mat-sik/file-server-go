package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"io"
	"os"
)

func HandleRequest(
	writer io.Writer,
	headerBuffer []byte,
	messageBuffer *bytes.Buffer,
	holder *message.Holder,
) error {
	defer messageBuffer.Reset()

	if err := transfer.SendMessage(writer, headerBuffer, messageBuffer, holder); err != nil {
		return err
	}
	return nil
}

func HandleResponse(
	ctx context.Context,
	reader io.Reader,
	buffer *bytes.Buffer,
) (message.Holder, error) {
	defer buffer.Reset()

	holder, err := transfer.ReceiveMessage(ctx, reader, buffer)
	if err != nil {
		return message.Holder{}, err
	}
	return holder, nil
}

func GetFileHandleResponse(
	ctx context.Context,
	reader io.Reader,
	buffer *bytes.Buffer,
	filename string,
) error {
	defer buffer.Reset()

	holder, err := transfer.ReceiveMessage(ctx, reader, buffer)
	if err != nil {
		return err
	}

	res, ok := holder.PayloadStruct.(*message.GetFileResponse)
	if !ok {
		return ErrUnexpectedResponse
	}
	if res.Status != 200 {
		return fmt.Errorf("error, status code: %d", res.Status)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	fileSize := res.Size
	if err = transfer.Stream(ctx, reader, file, buffer, fileSize); err != nil {
		return err
	}
	return nil
}

func PutFileHandleRequest(
	ctx context.Context,
	writer io.Writer,
	headerBuffer []byte,
	buffer *bytes.Buffer,
	filename string,
) error {
	defer buffer.Reset()

	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := int(fileInfo.Size())

	holder := message.NewPutFileRequestHolder(filename, fileSize)
	if err = HandleRequest(writer, headerBuffer, buffer, &holder); err != nil {
		return err
	}
	return transfer.Stream(ctx, file, writer, buffer, fileSize)
}

var ErrUnexpectedResponse = errors.New("unexpected response")
