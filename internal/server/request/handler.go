package request

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/connection"
)

func HandleGetFileRequest(req message.Request) (message.Response, error) {
	getFileReq := req.(*message.GetFileRequest)
	fileName := getFileReq.Filename

	return handleGetFileRequest(fileName)
}

func HandlePutFileRequest(ctx context.Context, connCtx connection.Context, req message.Request) (message.Response, error) {
	putFileReq := req.(*message.PutFileRequest)
	fileName := putFileReq.FileName
	fileSize := putFileReq.Size

	writer := connCtx.Conn
	buffer := connCtx.Buffer

	return handlePutFileRequest(ctx, writer, buffer, fileName, fileSize)
}

func HandleDeleteFileRequest(req message.Request) (message.Response, error) {
	deleteFileReq := req.(*message.DeleteFileRequest)
	fileName := deleteFileReq.FileName
	return handleDeleteFileRequest(fileName)
}
