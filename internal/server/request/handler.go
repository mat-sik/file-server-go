package request

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/envs"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/message/decorated"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
	"os"
	"path/filepath"
)

func HandleGetFileRequest(req message.GetFileRequest) decorated.GetFileResponse {
	res := message.GetFileResponse{}
	return decorated.New(res, req.FileName)
}

func HandlePutFileRequest(
	ctx context.Context,
	writer io.Writer,
	buffer *limited.Buffer,
	req message.PutFileRequest,
) (message.PutFileResponse, error) {
	defer buffer.Reset()

	path := filepath.Join(envs.ServerDBPath, req.FileName)
	file, err := os.Create(path)
	if err != nil {
		return message.PutFileResponse{}, err
	}

	if err = transfer.Stream(ctx, file, writer, buffer, req.Size); err != nil {
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
