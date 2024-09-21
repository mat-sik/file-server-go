package controller

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/server/router"
	"github.com/mat-sik/file-server-go/internal/server/service"
)

func GetFile(ctx context.Context, rs router.RequestState, req message.GetFileRequest) error {
	filename := req.Filename
	return service.GetFile(ctx, rs, filename)
}

func PutFile(ctx context.Context, rs router.RequestState, req message.PutFileRequest) (message.Holder, error) {
	filename := req.FileName
	fileSize := req.Size
	return service.PutFile(ctx, rs, filename, fileSize)
}

func DeleteFile(req message.DeleteFileRequest) (message.Holder, error) {
	filename := req.FileName
	return service.DeleteFile(filename)
}
