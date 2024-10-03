package request

import (
	"bytes"
	"context"
	"github.com/mat-sik/file-server-go/internal/envs"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"io"
	"os"
	"path/filepath"
)

func NewGetFileRequest(fileName string) (message.Request, error) {
	return &message.GetFileRequest{Filename: fileName}, nil
}

func NewPutFileRequest(fileName string) (message.Request, error) {
	path := filepath.Join(envs.ClientDBPath, fileName)
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := int(fileInfo.Size())

	req := message.PutFileRequest{FileName: fileName, Size: fileSize}
	return &StreamRequest{
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

func (req *StreamRequest) GetResponseType() message.RequestTypeName {
	return req.Request.GetRequestType()
}

func (req *StreamRequest) Stream(
	ctx context.Context,
	writer io.Writer,
	headerBuffer []byte,
	messageBuffer *bytes.Buffer,
) error {
	return transfer.StreamFromFile(ctx, writer, headerBuffer, messageBuffer, req)
}

func NewDeleteFileRequest(fileName string) (message.Request, error) {
	return &message.DeleteFileRequest{FileName: fileName}, nil
}
