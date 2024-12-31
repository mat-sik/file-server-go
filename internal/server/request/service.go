package request

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/envs"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/message/decorated"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
	"os"
	"path/filepath"
)

func handleGetFileRequest(fileName string) decorated.GetFileResponse {
	res := message.GetFileResponse{}
	return decorated.New(res, fileName)
}

func handlePutFileRequest(
	ctx context.Context,
	writer io.Writer,
	buffer *limited.Buffer,
	fileName string,
	fileSize int,
) (message.PutFileResponse, error) {
	defer buffer.Reset()

	path := filepath.Join(envs.ServerDBPath, fileName)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return message.PutFileResponse{}, err
	}

	if err = transfer.Stream(ctx, file, writer, buffer, fileSize); err != nil {
		return message.PutFileResponse{}, err
	}

	return message.PutFileResponse{Status: 200}, nil
}

func handleDeleteFileRequest(fileName string) (message.DeleteFileResponse, error) {
	err := os.Remove(fileName)
	if err != nil {
		return message.DeleteFileResponse{}, err
	}
	return message.DeleteFileResponse{Status: 200}, nil
}
