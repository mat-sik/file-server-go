package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/client/response"
	"github.com/mat-sik/file-server-go/internal/file"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/message/decorated"
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

	decorateRes := func(res message.GetFileResponse) decorated.GetFileResponse {
		req, ok := req.(*message.GetFileRequest)
		if !ok {
			panic(fmt.Sprintf("GetFileRequest expected, received: %v", req))
		}
		return decorated.GetFileResponse{GetFileResponse: res, FileName: req.FileName}
	}

	res, err := sh.receiveResponse(decorateRes)
	if err != nil {
		return err
	}

	return sh.handleResponse(ctx, res)
}

func (sh SessionHandler) deliverRequest(ctx context.Context, req message.Request) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req.GetType() {
	case message.PutFileRequestType:
		req := req.(*message.PutFileRequest)
		return sh.streamRequest(ctx, req)
	default:
		return sh.SendMessage(req)
	}
}

func (sh SessionHandler) streamRequest(ctx context.Context, req *message.PutFileRequest) error {
	defer sh.Buffer.Reset()

	f, err := os.Open(req.FileName)
	if err != nil {
		return err
	}
	defer file.Close(f)

	fileSize, err := file.GetSize(f)
	req.Size = fileSize

	if err = sh.SendMessage(req); err != nil {
		return err
	}
	return sh.StreamToNet(ctx, f, fileSize)
}

func (sh SessionHandler) receiveResponse(
	decorateRes func(fileResponse message.GetFileResponse) decorated.GetFileResponse,
) (message.Response, error) {
	m, err := sh.ReceiveMessage()
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

func (sh SessionHandler) handleResponse(ctx context.Context, res message.Response) error {
	defer sh.Buffer.Reset()

	switch res.GetType() {
	case message.GetFileResponseType:
		res := res.(decorated.GetFileResponse)
		return response.HandelGetFileResponse(ctx, sh.Session, res)
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
