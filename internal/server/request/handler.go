package request

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/envs"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/message/decorated"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"os"
	"path/filepath"
)

func HandleGetFileRequest(req message.GetFileRequest) decorated.GetFileResponse {
	res := message.GetFileResponse{}
	return decorated.GetFileResponse{GetFileResponse: res, FileName: req.FileName}
}

func HandlePutFileRequest(
	ctx context.Context,
	dispatcher netmsg.MessageDispatcher,
	req message.PutFileRequest,
) (message.PutFileResponse, error) {
	defer dispatcher.Buffer.Reset()

	path := filepath.Join(envs.ServerDBPath, req.FileName)
	file, err := os.Create(path)
	if err != nil {
		return message.PutFileResponse{}, err
	}

	if err = dispatcher.StreamFromNet(ctx, file, req.Size); err != nil {
		return message.PutFileResponse{}, err
	}

	return message.PutFileResponse{Status: 200}, nil
}

func HandleDeleteFileRequest(req message.DeleteFileRequest) (message.DeleteFileResponse, error) {
	err := os.Remove(req.FileName)
	if err != nil {
		return message.DeleteFileResponse{}, err
	}
	return message.DeleteFileResponse{Status: 200}, nil
}
