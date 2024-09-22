package router

import (
	"bytes"
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/server/controller"
	"github.com/mat-sik/file-server-go/internal/server/service"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/state"
	"io"
	"time"
)

func RouteRequest(ctx context.Context, s state.ConnectionState) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	m, err := transfer.ReceiveMessage(ctx, s.Conn, s.Buffer)
	if err != nil {
		return err
	}
	switch m.GetType() {
	case message.GetFileRequestType:
		if err = handleGetFile(ctx, s, m); err != nil {
			return err
		}
	case message.PutFileRequestType:
		if err = handlePutFile(ctx, s, m); err != nil {
			return err
		}
	case message.DeleteFileRequestType:
		if err = handleDeleteFile(s, m); err != nil {
			return err
		}
	default:
		return ErrUnexpectedRequestType
	}
	return nil
}

func handleGetFile(ctx context.Context, s state.ConnectionState, m message.Message) error {
	req := m.(*message.GetFileRequest)
	handleGetReq := func(req *message.GetFileRequest) (message.Response, error) {
		return controller.HandleGetFileRequest(s, req)
	}
	handleStreamGetRes := func(res message.Response) error {
		return handleStreamRes(ctx, s, res.(*service.StreamResponse))
	}
	if err := handleReq(req, handleGetReq, handleStreamGetRes); err != nil {
		return err
	}
	return nil
}

func handlePutFile(ctx context.Context, s state.ConnectionState, m message.Message) error {
	req := m.(*message.PutFileRequest)
	handlePutReq := func(req *message.PutFileRequest) (message.Response, error) {
		return controller.HandlePutFileRequest(ctx, s, req)
	}
	handlePutRes := func(res message.Response) error {
		return handleStreamRes(ctx, s, res.(*service.StreamResponse))
	}
	if err := handleReq(req, handlePutReq, handlePutRes); err != nil {
		return err
	}
	return nil
}

func handleDeleteFile(s state.ConnectionState, m message.Message) error {
	req := m.(*message.DeleteFileRequest)
	handleDelRes := func(res message.Response) error {
		writer := s.Conn
		headerBuffer := s.HeaderBuffer
		messageBuffer := s.Buffer
		return handleSendRes(writer, headerBuffer, messageBuffer, res)
	}
	if err := handleReq(req, controller.HandleDeleteFileRequest, handleDelRes); err != nil {
		return err
	}
	return nil
}

func handleReq[T message.Request](
	req T,
	handleReq func(T) (message.Response, error),
	handleRes func(message.Response) error,
) error {
	res, err := handleReq(req)
	if err != nil {
		return err
	}
	return handleRes(res)
}

func handleStreamRes(ctx context.Context, s state.ConnectionState, streamRes *service.StreamResponse) error {
	var writer io.Writer = s.Conn
	headerBuffer := s.HeaderBuffer
	buffer := s.Buffer

	res := streamRes.StructResponse
	if err := handleSendRes(writer, headerBuffer, buffer, res); err != nil {
		return err
	}

	reader := streamRes.StreamReader
	toTransfer := streamRes.ToTransfer
	if err := transfer.Stream(ctx, reader, writer, buffer, toTransfer); err != nil {
		return err
	}
	return nil
}

func handleSendRes(writer io.Writer, headerBuffer []byte, messageBuffer *bytes.Buffer, res message.Response) error {
	m := res.(message.Message)
	if err := transfer.SendMessage(writer, headerBuffer, messageBuffer, m); err != nil {
		return err
	}
	return nil
}

const timeForRequest = 5 * time.Second

var ErrUnexpectedRequestType = errors.New("unexpected request type")
