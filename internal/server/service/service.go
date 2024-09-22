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

func HandleGetFileRequest(messageBuffer *bytes.Buffer, filename string) (message.Response, error) {
	defer messageBuffer.Reset()

	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		return &message.GetFileResponse{Status: 404, Size: 0}, err
	}
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := int(fileInfo.Size())
	res := message.GetFileResponse{Status: 200, Size: fileSize}
	return &StreamResponse{StructResponse: &res, Reader: file, ToTransfer: fileSize}, nil
}

type StreamResponse struct {
	StructResponse message.Response
	io.Reader
	ToTransfer int
}

func (res *StreamResponse) GetResponseType() message.TypeName {
	return res.StructResponse.GetResponseType()
}

func HandlePutFileRequest(
	ctx context.Context,
	writer io.Writer,
	buffer *bytes.Buffer,
	filename string,
	fileSize int,
) (message.Response, error) {
	defer buffer.Reset()

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	if err = transfer.Stream(ctx, file, writer, buffer, fileSize); err != nil {
		return nil, err
	}

	return &message.PutFileResponse{Status: 200}, nil
}

func HandleDeleteFileRequest(filename string) (message.Response, error) {
	err := os.Remove(filename)
	if err != nil {
		return nil, err
	}
	return &message.DeleteFileResponse{Status: 200}, nil
}
