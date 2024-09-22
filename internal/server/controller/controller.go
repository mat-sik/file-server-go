package controller

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/server/service"
	"github.com/mat-sik/file-server-go/internal/transfer/state"
)

func HandleGetFileRequest(s state.ConnectionState, req *message.GetFileRequest) (message.Response, error) {
	messageBuffer := s.Buffer
	filename := req.Filename

	return service.HandleGetFileRequest(messageBuffer, filename)
}

func HandlePutFileRequest(ctx context.Context, s state.ConnectionState, req *message.PutFileRequest) (message.Response, error) {
	writer := s.Conn
	buffer := s.Buffer

	filename := req.FileName
	fileSize := req.Size

	return service.HandlePutFileRequest(ctx, writer, buffer, filename, fileSize)
}

func HandleDeleteFileRequest(req *message.DeleteFileRequest) (message.Response, error) {
	filename := req.FileName
	return service.HandleDeleteFileRequest(filename)
}
