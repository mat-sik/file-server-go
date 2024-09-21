package controller

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/server/router"
	"github.com/mat-sik/file-server-go/internal/server/service"
)

func GetFile(ctx context.Context, rs router.RequestState, req message.GetFileRequest) error {
	writer := rs.Conn
	headerBuffer := rs.HeaderBuffer
	messageBuffer := rs.Buffer

	filename := req.Filename

	return service.GetFile(ctx, writer, headerBuffer, messageBuffer, filename)
}

func PutFile(ctx context.Context, rs router.RequestState, req message.PutFileRequest) (message.Holder, error) {
	writer := rs.Conn
	buffer := rs.Buffer

	filename := req.FileName
	fileSize := req.Size

	return service.PutFile(ctx, writer, buffer, filename, fileSize)
}

func DeleteFile(req message.DeleteFileRequest) (message.Holder, error) {
	filename := req.FileName
	return service.DeleteFile(filename)
}
