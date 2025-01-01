package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/client/response"
	"github.com/mat-sik/file-server-go/internal/file"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/message/decorated"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"os"
	"time"
)

type ClientRouter struct {
	transfer.MessageDispatcher
}

func (clientRouter ClientRouter) HandleRequest(ctx context.Context, req message.Request) error {
	if err := clientRouter.deliverRequest(ctx, req); err != nil {
		return err
	}

	decorateRes := func(res message.GetFileResponse) decorated.GetFileResponse {
		req, ok := req.(*message.GetFileRequest)
		if !ok {
			panic(fmt.Sprintf("GetFileRequest expected, received: %v", req))
		}
		return decorated.GetFileResponse{GetFileResponse: res, FileName: req.FileName}
	}

	res, err := clientRouter.receiveResponse(decorateRes)
	if err != nil {
		return err
	}

	return clientRouter.handleResponse(ctx, res)
}

func (clientRouter ClientRouter) deliverRequest(ctx context.Context, req message.Request) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req.GetType() {
	case message.PutFileRequestType:
		req := req.(*message.PutFileRequest)
		return clientRouter.streamRequest(ctx, req)
	default:
		return clientRouter.SendMessage(req)
	}
}

func (clientRouter ClientRouter) streamRequest(ctx context.Context, req *message.PutFileRequest) error {
	defer clientRouter.Buffer.Reset()

	f, err := os.Open(req.FileName)
	if err != nil {
		return err
	}
	defer file.Close(f)

	fileSize, err := file.GetSize(f)
	req.Size = fileSize

	if err = clientRouter.SendMessage(req); err != nil {
		return err
	}
	return transfer.Stream(ctx, f, clientRouter, clientRouter.Buffer, fileSize)
}

func (clientRouter ClientRouter) receiveResponse(
	decorateRes func(fileResponse message.GetFileResponse) decorated.GetFileResponse,
) (message.Response, error) {
	m, err := clientRouter.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	res, ok := m.(message.Response)
	if !ok {
		return nil, errors.New("expected response, received different type")
	}

	if res.GetType() == message.GetFileResponseType {
		getFileResponse := res.(message.GetFileResponse)
		res = decorateRes(getFileResponse)
	}

	return res, nil
}

func (clientRouter ClientRouter) handleResponse(ctx context.Context, res message.Response) error {
	defer clientRouter.Buffer.Reset()

	switch res.GetType() {
	case message.GetFileResponseType:
		res := res.(decorated.GetFileResponse)
		return response.HandelGetFileResponse(ctx, clientRouter, clientRouter.Buffer, res)
	case message.PutFileResponseType:
		res := res.(message.PutFileResponse)
		response.HandlePutFileResponse(res)
	case message.DeleteFileResponseType:
		res := res.(message.DeleteFileResponse)
		response.HandleDeleteFileResponse(res)
	default:
		return errors.New("unexpected response type")
	}

	return nil
}

const timeForRequest = 5 * time.Second
