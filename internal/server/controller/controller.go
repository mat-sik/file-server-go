package controller

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/state"
)

func HandleGetFileRequest(req message.Request) (message.Response, error) {
	getFileReq := req.(*message.GetFileRequest)
	filename := getFileReq.Filename

	return handleGetFileRequest(filename)
}

func HandlePutFileRequest(ctx context.Context, s state.ConnectionState, req message.Request) (message.Response, error) {
	putFileReq := req.(*message.PutFileRequest)
	filename := putFileReq.FileName
	fileSize := putFileReq.Size

	writer := s.Conn
	buffer := s.Buffer

	return handlePutFileRequest(ctx, writer, buffer, filename, fileSize)
}

func HandleDeleteFileRequest(req message.Request) (message.Response, error) {
	deleteFileReq := req.(*message.DeleteFileRequest)
	filename := deleteFileReq.FileName
	return handleDeleteFileRequest(filename)
}
