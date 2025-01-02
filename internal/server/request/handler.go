package request

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/message/decorated"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"os"
)

func HandleGetFileRequest(req message.GetFileRequest) *decorated.GetFileResponse {
	res := message.GetFileResponse{}
	return &decorated.GetFileResponse{GetFileResponse: res, FileName: req.FileName}
}

func HandlePutFileRequest(
	ctx context.Context,
	session netmsg.Session,
	req message.PutFileRequest,
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

func HandleDeleteFileRequest(req message.DeleteFileRequest) (*message.DeleteFileResponse, error) {
	err := os.Remove(req.FileName)
	if err != nil {
		return nil, err
	}
	return &message.DeleteFileResponse{Status: 200}, nil
}
