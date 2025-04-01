package client

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"os"
	"time"
)

type sessionHandler struct {
	netmsg.Session
}

func (sh sessionHandler) handleRequest(ctx context.Context, req message.Request) (message.Response, error) {
	if err := sh.deliverRequest(ctx, req); err != nil {
		return nil, err
	}

	ctx = setValuesInContext(ctx, req)

	res, err := sh.receiveResponse()
	if err != nil {
		return nil, err
	}

	if err = sh.handleResponse(ctx, res); err != nil {
		return nil, err
	}

	return res, nil
}

func setValuesInContext(ctx context.Context, req message.Request) context.Context {
	if req, ok := req.(message.FilenameGetter); ok {
		ctx = contextWithFileName(ctx, req.GetFilename())
	}
	if req, ok := req.(message.GetFilenamesRequest); ok {
		ctx = contextWithPattern(ctx, req.MatchRegex)
	}
	return ctx
}

func (sh sessionHandler) deliverRequest(ctx context.Context, req message.Request) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req := req.(type) {
	case message.PutFileRequest:
		return sh.streamRequest(ctx, req)
	default:
		return sh.SendMessage(req)
	}
}

func (sh sessionHandler) streamRequest(ctx context.Context, req message.PutFileRequest) error {
	path := files.BuildClientFilePath(req.Filename)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer files.LoggedClose(file)

	fileSize, err := files.SizeOf(file)
	req.Size = fileSize

	if err = sh.SendMessage(req); err != nil {
		return err
	}
	return sh.StreamToNet(ctx, file, fileSize)
}

func (sh sessionHandler) receiveResponse() (message.Response, error) {
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

func (sh sessionHandler) handleResponse(ctx context.Context, res message.Response) error {
	switch res := res.(type) {
	case message.GetFileResponse:
		return sh.handleGetFileResponse(ctx, res)
	case message.PutFileResponse:
		handlePutFileResponse(ctx, res)
	case message.DeleteFileResponse:
		handleDeleteFileResponse(ctx, res)
	case message.GetFilenamesResponse:
		handleGetFilenamesResponse(ctx, res)
	default:
		return errors.New("unexpected response type")
	}

	return nil
}

func (sh sessionHandler) handleGetFileResponse(
	ctx context.Context,
	res message.GetFileResponse,
) error {
	filename := filenameFromContextOrPanic(ctx)
	return handelGetFileResponse(ctx, sh.Session, filename, res)
}

const timeForRequest = 5 * time.Second
