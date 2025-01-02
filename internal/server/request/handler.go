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
		return &GetFileResponse{GetFileResponse: message.NewGetFileResponse(http.StatusNotFound, 0)}, nil
	} else if err != nil {
		return nil, err
	}

	fileSize, err := files.GetSize(file)
	if err != nil {
		return nil, err
	}

	return &GetFileResponse{
		GetFileResponse: message.NewGetFileResponse(http.StatusOK, fileSize),
		File:            file,
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

	return message.NewPutFileResponse(http.StatusCreated), nil
}

func HandleDeleteFileRequest(req *message.DeleteFileRequest) (*message.DeleteFileResponse, error) {
	err := os.Remove(req.FileName)
	if errors.Is(err, os.ErrNotExist) {
		return message.NewDeleteFileResponse(http.StatusNotFound), nil
	}
	if err != nil {
		return nil, err
	}
	return message.NewDeleteFileResponse(http.StatusOK), nil
}
