package controller

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/server/service"
	"github.com/mat-sik/file-server-go/internal/transfer/state"
)

func GetFile(ctx context.Context, s state.ConnectionState, req message.GetFileRequest) error {
	writer := s.Conn
	headerBuffer := s.HeaderBuffer
	messageBuffer := s.Buffer

	filename := req.Filename

	return service.GetFile(ctx, writer, headerBuffer, messageBuffer, filename)
}

func PutFile(ctx context.Context, s state.ConnectionState, req message.PutFileRequest) (message.Holder, error) {
	writer := s.Conn
	buffer := s.Buffer

	filename := req.FileName
	fileSize := req.Size

	return service.PutFile(ctx, writer, buffer, filename, fileSize)
}

func DeleteFile(req message.DeleteFileRequest) (message.Holder, error) {
	filename := req.FileName
	return service.DeleteFile(filename)
}
