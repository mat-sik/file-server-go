package server

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/message/decorated"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"github.com/mat-sik/file-server-go/internal/server/request"
	"net/http"
	"os"
	"time"
)

type SessionHandler struct {
	netmsg.Session
}

func (sh SessionHandler) HandleRequest(ctx context.Context) error {
	req, err := sh.receiveRequest()
	if err != nil {
		return err
	}

	res, err := sh.routeRequest(ctx, req)
	if err != nil {
		return err
	}

	return sh.deliverResponse(ctx, res)
}

func (sh SessionHandler) receiveRequest() (message.Request, error) {
	msg, err := sh.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	req, ok := msg.(message.Request)
	if !ok {
		return nil, errors.New("expected request, received different type")
	}
	return req, nil
}

func (sh SessionHandler) routeRequest(ctx context.Context, req message.Request) (message.Response, error) {
	defer sh.Buffer.Reset()

	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req.GetType() {
	case message.GetFileRequestType:
		req := *req.(*message.GetFileRequest)
		return request.HandleGetFileRequest(req), nil
	case message.PutFileRequestType:
		req := *req.(*message.PutFileRequest)
		return request.HandlePutFileRequest(ctx, sh.Session, req)
	case message.DeleteFileRequestType:
		req := *req.(*message.DeleteFileRequest)
		return request.HandleDeleteFileRequest(req)
	default:
		return nil, errors.New("unexpected request type")
	}
}

func (sh SessionHandler) deliverResponse(ctx context.Context, res message.Response) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch res.GetType() {
	case message.GetFileResponseType:
		res := *res.(*decorated.GetFileResponse)
		return sh.sendGetFileResponse(ctx, res)
	default:
		return sh.SendMessage(res)
	}
}

func (sh SessionHandler) sendGetFileResponse(ctx context.Context, res decorated.GetFileResponse) error {
	path := files.GetServerDBPath(res.FileName)
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return sh.sendNotFoundResponse(res)
	} else if err != nil {
		return err
	}
	defer files.Close(file)

	return sh.streamFileResponse(ctx, file, res)
}

func (sh SessionHandler) streamFileResponse(
	ctx context.Context,
	file *os.File,
	res decorated.GetFileResponse,
) error {
	fileSize, err := files.GetSize(file)
	if err != nil {
		return err
	}
	res.GetFileResponse.Status = http.StatusOK
	res.Size = fileSize

	if err = sh.SendMessage(res.GetFileResponse); err != nil {
		return err
	}
	return sh.StreamToNet(ctx, file, res.Size)
}

func (sh SessionHandler) sendNotFoundResponse(res decorated.GetFileResponse) error {
	res.GetFileResponse.Status = http.StatusNotFound
	res.Size = 0

	return sh.SendMessage(res)
}

const timeForRequest = 5 * time.Second
