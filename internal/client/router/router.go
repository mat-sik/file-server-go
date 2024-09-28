package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/client/reshandler"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/state"
	"io"
	"time"
)

func HandleRequest(ctx context.Context, s state.ConnectionState, req message.Request) error {
	if err := deliverRequest(ctx, s, req); err != nil {
		return err
	}

	enrichFunc := func(res message.Response) message.Response {
		return enrichGetFileResponse(res, req)
	}

	res, err := receiveResponse(ctx, s, enrichFunc)
	if err != nil {
		return err
	}

	return handleResponse(ctx, s, res)
}

func receiveResponse(
	ctx context.Context,
	s state.ConnectionState,
	enrichFunc func(message.Response) message.Response,
) (message.Response, error) {
	var reader io.Reader = s.Conn
	buffer := s.Buffer
	m, err := transfer.ReceiveMessage(ctx, reader, buffer)
	if err != nil {
		return nil, err
	}

	res, ok := m.(message.Response)
	if !ok {
		return nil, ErrExpectedResponse
	}

	if res.GetResponseType() == message.GetFileResponseType {
		res = enrichFunc(res)
	}

	return res, nil
}

func enrichGetFileResponse(res message.Response, req message.Request) message.Response {
	getFileRequest, ok := req.(*message.GetFileRequest)
	if !ok {
		panic(fmt.Sprintf("GetFileRequest expected, received: %v", req))
	}
	filename := getFileRequest.Filename

	getFileResponse := res.(*message.GetFileResponse)

	return EnrichedGetFileResponse{
		Response: getFileResponse,
		Filename: filename,
	}
}

type EnrichedGetFileResponse struct {
	message.Response
	Filename string
}

func deliverRequest(ctx context.Context, s state.ConnectionState, req message.Request) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	switch req.GetRequestType() {
	case message.PutFileRequestType:
		return streamRequest(ctx, s, req)
	default:
		return sendRequest(s, req)
	}
}

func streamRequest(ctx context.Context, s state.ConnectionState, req message.Request) error {
	streamReq := req.(message.StreamableMessage)

	var writer io.Writer = s.Conn
	headerBuffer := s.HeaderBuffer
	messageBuffer := s.Buffer
	return streamReq.Stream(ctx, writer, headerBuffer, messageBuffer)
}

func sendRequest(s state.ConnectionState, req message.Request) error {
	m := req.(message.Message)

	var writer io.Writer = s.Conn
	headerBuffer := s.HeaderBuffer
	messageBuffer := s.Buffer
	return transfer.SendMessage(writer, headerBuffer, messageBuffer, m)
}

func handleResponse(ctx context.Context, s state.ConnectionState, res message.Response) error {
	buffer := s.Buffer
	defer buffer.Reset()

	switch res.GetResponseType() {
	case message.GetFileResponseType:
		return reshandler.HandelGetFileResponse(ctx, s, res)
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
