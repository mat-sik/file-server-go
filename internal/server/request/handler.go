package request

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"net/http"
	"os"
)

func HandleGetFileRequest(req message.GetFileRequest) (GetFileResponse, error) {
	path := files.BuildServerFilePath(req.FileName)
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return GetFileResponse{
			GetFileResponse: message.GetFileResponse{
				Status: http.StatusNotFound,
				Size:   0,
			},
		}, nil
	} else if err != nil {
		return GetFileResponse{}, err
	}

	fileSize, err := files.SizeOf(file)
	if err != nil {
		return GetFileResponse{}, err
	}

	return GetFileResponse{
		GetFileResponse: message.GetFileResponse{
			Status: http.StatusOK,
			Size:   fileSize,
		},
		File: file,
	}, nil
}

type GetFileResponse struct {
	message.GetFileResponse
	*os.File
}

func HandlePutFileRequest(
	ctx context.Context,
	session netmsg.Session,
	req message.PutFileRequest,
) (message.PutFileResponse, error) {
	path := files.BuildServerFilePath(req.FileName)
	file, err := os.Create(path)
	if err != nil {
		return message.PutFileResponse{}, err
	}

	if err = session.StreamFromNet(ctx, file, req.Size); err != nil {
		return message.PutFileResponse{}, err
	}

	return message.PutFileResponse{
		Status: http.StatusCreated,
	}, nil
}

func HandleDeleteFileRequest(req message.DeleteFileRequest) (message.DeleteFileResponse, error) {
	path := files.BuildServerFilePath(req.FileName)
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return message.DeleteFileResponse{
			Status: http.StatusNotFound,
		}, nil
	}
	if err != nil {
		return message.DeleteFileResponse{}, err
	}
	return message.DeleteFileResponse{
		Status: http.StatusOK,
	}, nil
}
