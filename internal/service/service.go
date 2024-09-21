package service

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/handler"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"os"
)

func GetFile(
	ctx context.Context,
	rs handler.RequestState,
	filename string,
) error {
	defer rs.Buffer.Reset()

	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		holder := message.NewGetFileResponseHolder(404, 0)
		return sendMessage(rs, holder)
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
	if err = sendMessage(rs, holder); err != nil {
		return err
	}

	writer := rs.Conn
	messageBuffer := rs.Buffer
	if err = transfer.Stream(ctx, file, writer, messageBuffer, fileSize); err != nil {
		return err
	}

	return nil
}

func sendMessage(rs handler.RequestState, holder message.Holder) error {
	rs.Buffer.Reset()
	defer rs.Buffer.Reset()

	writer := rs.Conn
	headerBuffer := rs.HeaderBuffer
	messageBuffer := rs.Buffer
	return transfer.SendMessage(writer, headerBuffer, messageBuffer, &holder)
}

func PutFile(ctx context.Context, rs handler.RequestState, filename string, fileSize int) (message.Holder, error) {
	defer rs.Buffer.Reset()
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return message.Holder{}, err
	}

	writer := rs.Conn
	buffer := rs.Buffer
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
