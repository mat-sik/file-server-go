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
	"github.com/mat-sik/file-server-go/internal/transfer/connection"
	"io"
	"os"
	"time"
)

func HandleRequest(ctx context.Context, connCtx connection.Context, req message.Request) error {
	if err := deliverRequest(ctx, connCtx, req); err != nil {
		return err
	}

	decorateRes := func(res message.GetFileResponse) decorated.GetFileResponse {
		req, ok := req.(*message.GetFileRequest)
		if !ok {
			panic(fmt.Sprintf("GetFileRequest expected, received: %v", req))
		}
		return decorated.New(res, req.FileName)
	}

	res, err := receiveResponse(connCtx, decorateRes)
	if err != nil {
		return err
	}

	return handleResponse(ctx, connCtx, res)
}

func deliverRequest(ctx context.Context, connCtx connection.Context, req message.Request) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req.GetType() {
	case message.PutFileRequestType:
		req := req.(*message.PutFileRequest)
		return streamRequest(ctx, connCtx, req)
	default:
		return sendRequest(connCtx, req)
	}
}

func streamRequest(ctx context.Context, connCtx connection.Context, req *message.PutFileRequest) error {
	var writer io.Writer = connCtx.Conn
	headerBuffer := connCtx.HeaderBuffer
	messageBuffer := connCtx.Buffer

	defer messageBuffer.Reset()

	f, err := os.Open(req.FileName)
	if err != nil {
		return err
	}
	defer file.Close(f)

	fileSize, err := file.GetSize(f)
	req.Size = fileSize

	if err = transfer.SendMessage(writer, headerBuffer, messageBuffer, req); err != nil {
		return err
	}
	return transfer.Stream(ctx, f, writer, messageBuffer, fileSize)
}

func sendRequest(connCtx connection.Context, req message.Request) error {
	var writer io.Writer = connCtx.Conn
	headerBuffer := connCtx.HeaderBuffer
	messageBuffer := connCtx.Buffer
	return transfer.SendMessage(writer, headerBuffer, messageBuffer, req)
}

func receiveResponse(
	s connection.Context,
	decorateRes func(fileResponse message.GetFileResponse) decorated.GetFileResponse,
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

	if res.GetType() == message.GetFileResponseType {
		getFileResponse := res.(message.GetFileResponse)
		res = decorateRes(getFileResponse)
	}

	return res, nil
}

func handleResponse(ctx context.Context, connCtx connection.Context, res message.Response) error {
	buffer := connCtx.Buffer
	defer buffer.Reset()

	switch res.GetType() {
	case message.GetFileResponseType:
		res := res.(decorated.GetFileResponse)
		return response.HandelGetFileResponse(ctx, connCtx, res)
	case message.PutFileResponseType:
		res := res.(message.PutFileResponse)
		response.HandlePutFileResponse(res)
	case message.DeleteFileResponseType:
		res := res.(message.DeleteFileResponse)
		response.HandleDeleteFileResponse(res)
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
