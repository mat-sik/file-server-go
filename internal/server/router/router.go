package router

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/server/request"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/connection"
	"io"
	"time"
)

func HandleRequest(ctx context.Context, connCtx connection.Context) error {
	req, err := receiveRequest(connCtx)
	if err != nil {
		return err
	}

	res, err := routeRequest(ctx, connCtx, req)
	if err != nil {
		return err
	}

	return deliverResponse(ctx, connCtx, res)
}

func receiveRequest(connCtx connection.Context) (message.Request, error) {
	var reader io.Reader = connCtx.Conn
	buffer := connCtx.Buffer
	m, err := transfer.ReceiveMessage(reader, buffer)
	if err != nil {
		return nil, err
	}

	req, ok := m.(message.Request)
	if !ok {
		return nil, ErrExpectedRequest
	}
	return req, nil
}

func routeRequest(ctx context.Context, connCtx connection.Context, req message.Request) (message.Response, error) {
	buffer := connCtx.Buffer
	defer buffer.Reset()

	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req.GetType() {
	case message.GetFileRequestType:
		req := req.(message.GetFileRequest)
		return request.HandleGetFileRequest(req)
	case message.PutFileRequestType:
		req := req.(message.PutFileRequest)
		return request.HandlePutFileRequest(ctx, connCtx, req)
	case message.DeleteFileRequestType:
		req := req.(message.DeleteFileRequest)
		return request.HandleDeleteFileRequest(req)
	default:
		return nil, ErrUnexpectedRequestType
	}
}

func deliverResponse(ctx context.Context, connCtx connection.Context, res message.Response) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch res.GetType() {
	case message.GetFileResponseType:
		res := res.(message.StreamableMessage)
		return streamResponse(ctx, connCtx, res)
	default:
		return sendResponse(connCtx, res)
	}
}

func streamResponse(ctx context.Context, connCtx connection.Context, res message.StreamableMessage) error {
	var writer io.Writer = connCtx.Conn
	headerBuffer := connCtx.HeaderBuffer
	messageBuffer := connCtx.Buffer
	return res.Stream(ctx, writer, headerBuffer, messageBuffer)
}

func sendResponse(connCtx connection.Context, res message.Response) error {
	var writer io.Writer = connCtx.Conn
	headerBuffer := connCtx.HeaderBuffer
	messageBuffer := connCtx.Buffer
	return transfer.SendMessage(writer, headerBuffer, messageBuffer, res)
}

const timeForRequest = 5 * time.Second

var (
	ErrUnexpectedRequestType = errors.New("unexpected request type")
	ErrExpectedRequest       = errors.New("expected request, received different type")
)
