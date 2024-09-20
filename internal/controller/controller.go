package controller

import (
	"bytes"
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/service"
	"net"
)

type RequestState struct {
	Conn   net.Conn
	Buffer *bytes.Buffer
}

func getFile(ctx context.Context, rs RequestState, req message.GetFileRequest) (message.Holder, error) {
	filename := req.Filename
	return service.GetFile(ctx, rs, filename)
}

func putFile(ctx context.Context, rs RequestState, req message.PutFileRequest) (message.Holder, error) {
	filename := req.FileName
	fileSize := req.Size
	return service.PutFile(ctx, rs, filename, fileSize)
}
