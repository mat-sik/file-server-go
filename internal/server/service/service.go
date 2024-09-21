package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"io"
	"os"
)

func GetFile(
	ctx context.Context,
	writer io.Writer,
	headerBuffer []byte,
	messageBuffer *bytes.Buffer,
	filename string,
) error {
	defer messageBuffer.Reset()

	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		holder := message.NewGetFileResponseHolder(404, 0)
		return sendMessage(writer, headerBuffer, messageBuffer, holder)
	}
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := int(fileInfo.Size())
	holder := message.NewGetFileResponseHolder(200, fileSize)
	if err = sendMessage(writer, headerBuffer, messageBuffer, holder); err != nil {
		return err
	}

	if err = transfer.Stream(ctx, file, writer, messageBuffer, fileSize); err != nil {
		return err
	}

	return nil
}

func sendMessage(writer io.Writer, headerBuffer []byte, messageBuffer *bytes.Buffer, holder message.Holder) error {
	messageBuffer.Reset()
	defer messageBuffer.Reset()
	return transfer.SendMessage(writer, headerBuffer, messageBuffer, &holder)
}

func PutFile(
	ctx context.Context,
	writer io.Writer,
	buffer *bytes.Buffer,
	filename string,
	fileSize int,
) (message.Holder, error) {
	defer buffer.Reset()

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return message.Holder{}, err
	}

	if err = transfer.Stream(ctx, file, writer, buffer, fileSize); err != nil {
		return message.Holder{}, err
	}

	return message.NewPutFileResponseHolder(200), nil
}

func DeleteFile(filename string) (message.Holder, error) {
	err := os.Remove(filename)
	if err != nil {
		return message.Holder{}, err
	}
	return message.NewDeleteFileResponseHolder(200), nil
}
