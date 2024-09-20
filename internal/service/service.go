package service

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/controller"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"os"
)

func GetFile(
	ctx context.Context,
	rs controller.RequestState,
	filename string,
) (message.Holder, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		return message.NewGetFileResponseHolder(404, 0), nil
	}
	if err != nil {
		return message.Holder{}, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return message.Holder{}, err
	}

	fileSize := int(fileInfo.Size())
	writer := rs.Conn
	buffer := rs.Buffer
	if err = transfer.Stream(ctx, file, writer, buffer, fileSize); err != nil {
		return message.Holder{}, err
	}

	return message.NewGetFileResponseHolder(200, fileSize), nil
}

func PutFile(ctx context.Context, rs controller.RequestState, filename string, fileSize int) (message.Holder, error) {
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

func DeleteFile(ctx context.Context, filename string) (message.Holder, error) {
	select {
	case <-ctx.Done():
	default:
		return message.Holder{}, ctx.Err()
	}
	err := os.Remove(filename)
	if err != nil {
		return message.Holder{}, err
	}
	return message.NewDeleteFileResponseHolder(200), nil
}
