package router

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/client/request/enricher"
	"github.com/mat-sik/file-server-go/internal/client/reshandler"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/conncontext"
	"io"
	"time"
)

func HandleRequest(ctx context.Context, connCtx conncontext.ConnectionContext, req message.Request) error {
	if err := deliverRequest(ctx, connCtx, req); err != nil {
		return err
	}

	enrichRes := func(res message.Response) message.Response {
		return enricher.EnrichGetFileResponse(res, req)
	}

	res, err := receiveResponse(connCtx, enrichRes)
	if err != nil {
		return err
	}

	return handleResponse(ctx, connCtx, res)
}

func receiveResponse(
	s conncontext.ConnectionContext,
	enrichRes func(message.Response) message.Response,
) (message.Response, error) {
	var reader io.Reader = s.Conn
	buffer := s.Buffer
	m, err := transfer.ReceiveMessage(reader, buffer)
	if err != nil {
		return nil, err
	}

	res, ok := m.(message.Response)
	if !ok {
		return nil, ErrExpectedResponse
	}

	if res.GetResponseType() == message.GetFileResponseType {
		res = enrichRes(res)
	}

	return res, nil
}

func deliverRequest(ctx context.Context, connCtx conncontext.ConnectionContext, req message.Request) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req.GetRequestType() {
	case message.PutFileRequestType:
		return streamRequest(ctx, connCtx, req)
	default:
		return sendRequest(connCtx, req)
	}
}

func streamRequest(ctx context.Context, connCtx conncontext.ConnectionContext, req message.Request) error {
	streamReq := req.(message.StreamableMessage)

	var writer io.Writer = connCtx.Conn
	headerBuffer := connCtx.HeaderBuffer
	messageBuffer := connCtx.Buffer
	return streamReq.Stream(ctx, writer, headerBuffer, messageBuffer)
}

func sendRequest(connCtx conncontext.ConnectionContext, req message.Request) error {
	m := req.(message.Message)

	var writer io.Writer = connCtx.Conn
	headerBuffer := connCtx.HeaderBuffer
	messageBuffer := connCtx.Buffer
	return transfer.SendMessage(writer, headerBuffer, messageBuffer, m)
}

func handleResponse(ctx context.Context, connCtx conncontext.ConnectionContext, res message.Response) error {
	buffer := connCtx.Buffer
	defer buffer.Reset()

	switch res.GetResponseType() {
	case message.GetFileResponseType:
		return reshandler.HandelGetFileResponse(ctx, connCtx, res)
	case message.PutFileResponseType:
		reshandler.HandlePutFileResponse(res)
	case message.DeleteFileResponseType:
		reshandler.HandleDeleteFileResponse(res)
	default:
		return ErrUnexpectedResponseType
	}

	return nil
}

const timeForRequest = 5 * time.Second

var (
	ErrExpectedResponse       = errors.New("expected response, received different type")
	ErrUnexpectedResponseType = errors.New("unexpected response type")
)
