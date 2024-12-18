package request

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/connection"
)

func HandleGetFileRequest(req message.GetFileRequest) (message.Response, error) {
	return handleGetFileRequest(req.Filename)
}

func HandlePutFileRequest(ctx context.Context, connCtx connection.Context, req message.PutFileRequest) (message.Response, error) {
	writer := connCtx.Conn
	buffer := connCtx.Buffer
	return handlePutFileRequest(ctx, writer, buffer, req.FileName, req.Size)
}

func HandleDeleteFileRequest(req message.DeleteFileRequest) (message.Response, error) {
	return handleDeleteFileRequest(req.FileName)
}
