package service

import (
	"github.com/mat-sik/file-server-go/internal/message"
	"io"
	"os"
)

func HandleGetFileRequest(filename string) (message.Request, error) {
	return &message.GetFileRequest{Filename: filename}, nil
}

func HandlePutFileRequest(filename string) (message.Request, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := int(fileInfo.Size())

	req := message.PutFileRequest{FileName: filename, Size: fileSize}
	return &StreamRequest{
		StructRequest: &req,
		Reader:        file,
		ToTransfer:    fileSize,
	}, nil
}

func HandleDeleteFileRequest(filename string) (message.Request, error) {
	return &message.DeleteFileRequest{FileName: filename}, nil
}

type StreamRequest struct {
	StructRequest message.Request
	io.Reader
	ToTransfer int
}

func (req *StreamRequest) GetRequestType() message.TypeName {
	return req.StructRequest.GetRequestType()
}
