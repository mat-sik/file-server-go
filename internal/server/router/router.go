package router

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/server/controller"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/state"
	"io"
	"time"
)

func HandleRequest(ctx context.Context, s state.ConnectionState) error {
	var reader io.Reader = s.Conn
	buffer := s.Buffer
	m, err := transfer.ReceiveMessage(ctx, reader, buffer)
	if err != nil {
		return err
	}

	req, ok := m.(message.Request)
	if !ok {
		return ErrExpectedRequest
	}

	res, err := RouteRequest(ctx, s, req)
	if err != nil {
		return err
	}

	return DeliverResponse(ctx, s, res)
}

func RouteRequest(ctx context.Context, s state.ConnectionState, req message.Request) (message.Response, error) {
	buffer := s.Buffer
	defer buffer.Reset()

	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req.GetRequestType() {
	case message.GetFileRequestType:
		return controller.HandleGetFileRequest(req)
	case message.PutFileRequestType:
		return controller.HandlePutFileRequest(ctx, s, req)
	case message.DeleteFileRequestType:
		return controller.HandleDeleteFileRequest(req)
	default:
		return nil, ErrUnexpectedRequestType
	}
}

func DeliverResponse(ctx context.Context, s state.ConnectionState, res message.Response) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch res.GetResponseType() {
	case message.GetFileResponseType:
		return streamResponse(ctx, s, res)
	default:
		return sendResponse(s, res)
	}
}

func streamResponse(ctx context.Context, s state.ConnectionState, res message.Response) error {
	streamRes := res.(message.StreamableMessage)

	var writer io.Writer = s.Conn
	headerBuffer := s.HeaderBuffer
	messageBuffer := s.Buffer
	return streamRes.Stream(ctx, writer, headerBuffer, messageBuffer)
}

func sendResponse(s state.ConnectionState, res message.Response) error {
	m := res.(message.Message)

	var writer io.Writer = s.Conn
	headerBuffer := s.HeaderBuffer
	messageBuffer := s.Buffer
	return transfer.SendMessage(writer, headerBuffer, messageBuffer, m)
}

const timeForRequest = 5 * time.Second

var (
	ErrUnexpectedRequestType = errors.New("unexpected request type")
	ErrExpectedRequest       = errors.New("expected request, received different type")
)
