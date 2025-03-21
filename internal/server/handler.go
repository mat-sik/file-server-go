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

type sessionHandler struct {
	netmsg.Session
	request.Handler
}

func (sh sessionHandler) handleRequest(ctx context.Context) error {
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

func (sh sessionHandler) receiveRequest() (message.Request, error) {
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

func (sh sessionHandler) routeRequest(ctx context.Context, req message.Request) (message.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req := req.(type) {
	case message.GetFileRequest:
		return sh.HandleGetFileRequest(req)
	case message.PutFileRequest:
		return sh.HandlePutFileRequest(ctx, sh.Session, req)
	case message.DeleteFileRequest:
		return sh.HandleDeleteFileRequest(req)
	case message.GetFilenamesRequest:
		return sh.HandleGetFilenamesRequest(req)
	default:
		return nil, errors.New("unexpected request type")
	}
}

func (sh sessionHandler) deliverResponse(ctx context.Context, res message.Response) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch res := res.(type) {
	case request.GetFileResponse:
		return sh.streamFileResponse(ctx, res)
	default:
		return sh.SendMessage(res)
	}
}

func (sh sessionHandler) streamFileResponse(ctx context.Context, res request.GetFileResponse) error {
	if res.Status != http.StatusOK {
		return sh.SendMessage(res.GetFileResponse)
	}

	defer files.LoggedClose(&res.ReadLockedFile)
	if err := sh.SendMessage(res.GetFileResponse); err != nil {
		return err
	}
	return sh.StreamToNet(ctx, res.ReadLockedFile, res.Size)
}

const timeForRequest = 5 * time.Second
