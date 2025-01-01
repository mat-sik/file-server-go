package router

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/file"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/message/decorated"
	"github.com/mat-sik/file-server-go/internal/server/request"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"net/http"
	"os"
	"time"
)

type ServerRouter struct {
	transfer.MessageDispatcher
}

func (serverRouter ServerRouter) HandleRequest(ctx context.Context) error {
	req, err := serverRouter.receiveRequest()
	if err != nil {
		return err
	}

	res, err := serverRouter.routeRequest(ctx, req)
	if err != nil {
		return err
	}

	return serverRouter.deliverResponse(ctx, res)
}

func (serverRouter ServerRouter) receiveRequest() (message.Request, error) {
	m, err := serverRouter.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	req, ok := m.(message.Request)
	if !ok {
		return nil, errors.New("expected request, received different type")
	}
	return req, nil
}

func (serverRouter ServerRouter) routeRequest(ctx context.Context, req message.Request) (message.Response, error) {
	defer serverRouter.Buffer.Reset()

	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req.GetType() {
	case message.GetFileRequestType:
		req := req.(message.GetFileRequest)
		return request.HandleGetFileRequest(req), nil
	case message.PutFileRequestType:
		req := req.(message.PutFileRequest)
		return request.HandlePutFileRequest(ctx, serverRouter, serverRouter.Buffer, req)
	case message.DeleteFileRequestType:
		req := req.(message.DeleteFileRequest)
		return request.HandleDeleteFileRequest(req)
	default:
		return nil, errors.New("unexpected request type")
	}
}

func (serverRouter ServerRouter) deliverResponse(ctx context.Context, res message.Response) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch res.GetType() {
	case message.GetFileResponseType:
		res := res.(decorated.GetFileResponse)
		return serverRouter.sendGetFileResponse(ctx, res)
	default:
		return serverRouter.SendMessage(res)
	}
}

func (serverRouter ServerRouter) sendGetFileResponse(ctx context.Context, res decorated.GetFileResponse) error {
	f, err := os.Open(res.FileName)
	if errors.Is(err, os.ErrNotExist) {
		return serverRouter.sendNotFoundResponse(res)
	} else if err != nil {
		return err
	}
	defer file.Close(f)

	return serverRouter.streamFileResponse(ctx, f, res)
}

func (serverRouter ServerRouter) streamFileResponse(
	ctx context.Context,
	f *os.File,
	res decorated.GetFileResponse,
) error {
	fileSize, err := file.GetSize(f)
	if err != nil {
		return err
	}
	res.Size = fileSize

	if err = serverRouter.SendMessage(res.GetFileResponse); err != nil {
		return err
	}
	return transfer.Stream(ctx, f, serverRouter, serverRouter.Buffer, res.Size)
}

func (serverRouter ServerRouter) sendNotFoundResponse(res decorated.GetFileResponse) error {
	res.GetFileResponse.Status = http.StatusNotFound
	res.Size = 0

	return serverRouter.SendMessage(res)
}

const timeForRequest = 5 * time.Second
