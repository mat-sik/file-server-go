package request

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/envs"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
	"os"
	"path/filepath"
)

func handleGetFileRequest(fileName string) (StreamResponse, error) {
	path := filepath.Join(envs.ServerDBPath, fileName)
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		return StreamResponse{GetFileResponse: message.GetFileResponse{Status: 404, Size: 0}}, err
	}
	if err != nil {
		return StreamResponse{}, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return StreamResponse{}, err
	}

	fileSize := int(fileInfo.Size())
	res := message.GetFileResponse{Status: 200, Size: fileSize}
	return StreamResponse{GetFileResponse: res, File: file, ToTransfer: fileSize}, nil
}

type StreamResponse struct {
	message.GetFileResponse
	*os.File
	ToTransfer int
}

func (res *StreamResponse) GetMessage() message.Message {
	return res.Response.(message.Message)
}

func (res *StreamResponse) GetFile() *os.File {
	return res.File
}

func (res *StreamResponse) GetToTransfer() int {
	return res.ToTransfer
}

func (res *StreamResponse) Stream(
	ctx context.Context,
	writer io.Writer,
	headerBuffer []byte,
	messageBuffer *limited.Buffer,
) error {
	return transfer.StreamFromFile(ctx, writer, headerBuffer, messageBuffer, res)
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
