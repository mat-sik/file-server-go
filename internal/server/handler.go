package server

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"github.com/mat-sik/file-server-go/internal/server/request"
	"net/http"
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
		return request.HandleGetFileRequest(req)
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
		res := *res.(*request.GetFileResponse)
		return sh.streamFileResponse(ctx, res)
	default:
		return sh.SendMessage(res)
	}
}

func (sh SessionHandler) streamFileResponse(
	ctx context.Context,
	res request.GetFileResponse,
) error {
	if res.Status != http.StatusOK {
		return sh.SendMessage(res)
	}

	defer files.Close(res.File)
	if err := sh.SendMessage(res.GetFileResponse); err != nil {
		return err
	}
	return sh.StreamToNet(ctx, res.File, res.Size)
}

const timeForRequest = 5 * time.Second
