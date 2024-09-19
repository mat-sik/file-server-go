package controller

import (
	"bytes"
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"net"
	"os"
)

type RequestState struct {
	conn   net.Conn
	buffer *bytes.Buffer
}

func (rs RequestState) getFile(ctx context.Context, req message.GetFileRequest) (*message.GetFileResponse, error) {
	filename := req.Filename

	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		return &message.GetFileResponse{Status: 404, Size: 0}, nil
	}
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := int(fileInfo.Size())
	if err = transfer.Stream(ctx, file, rs.conn, rs.buffer, fileSize); err != nil {
		return nil, err
	}

	return &message.GetFileResponse{Status: 200, Size: fileSize}, nil
}
