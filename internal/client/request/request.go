package request

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/envs"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
	"os"
	"path/filepath"
)

func NewGetFileRequest(fileName string) (message.GetFileRequest, error) {
	return message.GetFileRequest{Filename: fileName}, nil
}

func NewPutFileRequest(fileName string) (StreamRequest, error) {
	path := filepath.Join(envs.ClientDBPath, fileName)
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return StreamRequest{}, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return StreamRequest{}, err
	}

	fileSize := int(fileInfo.Size())

	req := message.PutFileRequest{FileName: fileName, Size: fileSize}
	return StreamRequest{
		Request:    &req,
		File:       file,
		ToTransfer: fileSize,
	}, nil
}

type StreamRequest struct {
	message.Request
	*os.File
	ToTransfer int
}

func (req *StreamRequest) GetMessage() message.Message {
	return req.Request.(message.Message)
}

func (req *StreamRequest) GetFile() *os.File {
	return req.File
}

func (req *StreamRequest) GetToTransfer() int {
	return req.ToTransfer
}

func (req *StreamRequest) Stream(
	ctx context.Context,
	writer io.Writer,
	headerBuffer []byte,
	messageBuffer *limited.Buffer,
) error {
	return transfer.StreamFromFile(ctx, writer, headerBuffer, messageBuffer, req)
}

func NewDeleteFileRequest(fileName string) (message.DeleteFileRequest, error) {
	return message.DeleteFileRequest{FileName: fileName}, nil
}
