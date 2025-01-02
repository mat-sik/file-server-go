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

func HandleGetFileRequest(req *message.GetFileRequest) (*GetFileResponse, error) {
	path := files.GetServerDBPath(req.FileName)
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return &GetFileResponse{GetFileResponse: &message.GetFileResponse{Status: http.StatusNotFound}}, nil
	} else if err != nil {
		return nil, err
	}

	fileSize, err := files.GetSize(file)
	if err != nil {
		return nil, err
	}

	return &GetFileResponse{
		GetFileResponse: &message.GetFileResponse{
			Status: http.StatusOK,
			Size:   fileSize,
		},
		File: file,
	}, nil
}

type GetFileResponse struct {
	*message.GetFileResponse
	*os.File
}

func HandlePutFileRequest(
	ctx context.Context,
	session netmsg.Session,
	req *message.PutFileRequest,
) (*message.PutFileResponse, error) {
	defer session.Buffer.Reset()

	path := files.GetServerDBPath(req.FileName)
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	if err = session.StreamFromNet(ctx, file, req.Size); err != nil {
		return nil, err
	}

	return &message.PutFileResponse{Status: 200}, nil
}

func HandleDeleteFileRequest(req *message.DeleteFileRequest) (*message.DeleteFileResponse, error) {
	err := os.Remove(req.FileName)
	if err != nil {
		return nil, err
	}
	return &message.DeleteFileResponse{Status: 200}, nil
}
