package router

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/file"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/message/decorated"
	"github.com/mat-sik/file-server-go/internal/server/request"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/connection"
	"io"
	"net/http"
	"os"
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
		return request.HandleGetFileRequest(req), nil
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
		res := res.(decorated.GetFileResponse)
		return sendGetFileResponse(ctx, connCtx, res)
	default:
		return sendResponse(connCtx, res)
	}
}

func sendGetFileResponse(ctx context.Context, connCtx connection.Context, res decorated.GetFileResponse) error {
	f, err := os.Open(res.FileName)
	if errors.Is(err, os.ErrNotExist) {
		return sendNotFoundResponse(connCtx, res)
	} else if err != nil {
		return err
	}
	defer file.Close(f)

	return streamFileResponse(ctx, connCtx, f, res)
}

func streamFileResponse(
	ctx context.Context,
	connCtx connection.Context,
	f *os.File,
	res decorated.GetFileResponse,
) error {
	var writer io.Writer = connCtx.Conn
	headerBuffer := connCtx.HeaderBuffer
	messageBuffer := connCtx.Buffer

	fileSize, err := file.GetSize(f)
	if err != nil {
		return err
	}
	res.Size = fileSize

	if err = transfer.SendMessage(writer, headerBuffer, messageBuffer, res.GetFileResponse); err != nil {
		return err
	}
	return transfer.Stream(ctx, f, writer, messageBuffer, res.Size)
}

func sendNotFoundResponse(connCtx connection.Context, res decorated.GetFileResponse) error {
	res.GetFileResponse.Status = http.StatusNotFound
	res.Size = 0

	return sendResponse(connCtx, res)
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
