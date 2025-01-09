package client

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/client/response"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"os"
	"time"
)

type SessionHandler struct {
	netmsg.Session
}

func (sh SessionHandler) HandleRequest(ctx context.Context, req message.Request) error {
	if err := sh.deliverRequest(ctx, req); err != nil {
		return err
	}

	if req, ok := req.(message.GetFileRequest); ok {
		ctx = contextWithFileName(ctx, req.FileName)
	}

	res, err := sh.receiveResponse()
	if err != nil {
		return err
	}

	return sh.handleResponse(ctx, res)
}

func (sh SessionHandler) deliverRequest(ctx context.Context, req message.Request) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req := req.(type) {
	case message.PutFileRequest:
		return sh.streamRequest(ctx, req)
	default:
		return sh.SendMessage(req)
	}
}

func (sh SessionHandler) streamRequest(ctx context.Context, req message.PutFileRequest) error {
	path := files.BuildClientFilePath(req.FileName)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer files.Close(file)

	fileSize, err := files.SizeOf(file)
	req.Size = fileSize

	if err = sh.SendMessage(req); err != nil {
		return err
	}
	return sh.StreamToNet(ctx, file, fileSize)
}

func (sh SessionHandler) receiveResponse() (message.Response, error) {
	msg, err := sh.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	res, ok := msg.(message.Response)
	if !ok {
		return nil, errors.New("expected response, received different type")
	}

	return res, nil
}

func (sh SessionHandler) handleResponse(ctx context.Context, res message.Response) error {
	defer sh.Buffer.Reset()

	switch res := res.(type) {
	case message.GetFileResponse:
		return sh.handleGetFileResponse(ctx, res)
	case message.PutFileResponse:
		response.HandlePutFileResponse(res)
	case message.DeleteFileResponse:
		response.HandleDeleteFileResponse(res)
	default:
		return errors.New("unexpected response type")
	}

	return nil
}

func (sh SessionHandler) handleGetFileResponse(
	ctx context.Context,
	res message.GetFileResponse,
) error {
	fileName, ok := fileNameFromContext(ctx)
	if !ok {
		return errors.New("file name not found in the context")
	}
	return response.HandelGetFileResponse(ctx, sh.Session, fileName, res)
}

func contextWithFileName(ctx context.Context, fileName string) context.Context {
	return context.WithValue(ctx, fileNameKey{}, fileName)
}

func fileNameFromContext(ctx context.Context) (string, bool) {
	res, ok := ctx.Value(fileNameKey{}).(string)
	return res, ok
}

type fileNameKey struct{}

const timeForRequest = 5 * time.Second
