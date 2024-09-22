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
	return &StreamResponse{Response: &res, File: file, ToTransfer: fileSize}, nil
}

type StreamResponse struct {
	message.Response
	*os.File
	ToTransfer int
}

func (res *StreamResponse) GetResponseType() message.ResponseTypeName {
	return res.Response.GetResponseType()
}

func (res *StreamResponse) Stream(
	ctx context.Context,
	writer io.Writer,
	headerBuffer []byte,
	messageBuffer *bytes.Buffer,
) error {
	file := res.File
	defer safeFileClose(file)

	m := res.Response.(message.Message)
	if err := transfer.SendMessage(writer, headerBuffer, messageBuffer, m); err != nil {
		return err
	}

	toTransfer := res.ToTransfer
	return transfer.Stream(ctx, file, writer, messageBuffer, toTransfer)
}

func safeFileClose(f *os.File) {
	if err := f.Close(); err != nil {
		panic(err)
	}
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
