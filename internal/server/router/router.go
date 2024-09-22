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
	handleReqFunc := func(req *message.GetFileRequest) (message.Response, error) {
		return controller.HandleGetFileRequest(s, req)
	}
	deliverResFunc := func(res message.Response) error {
		return deliverStreamRes(ctx, s, res.(*service.StreamResponse))
	}
	if err := handleReq(req, handleReqFunc, deliverResFunc); err != nil {
		return err
	}
	return nil
}

func handlePutFile(ctx context.Context, s state.ConnectionState, m message.Message) error {
	req := m.(*message.PutFileRequest)
	handleReqFunc := func(req *message.PutFileRequest) (message.Response, error) {
		return controller.HandlePutFileRequest(ctx, s, req)
	}
	deliverResFunc := func(res message.Response) error {
		return deliverStreamRes(ctx, s, res.(*service.StreamResponse))
	}
	if err := handleReq(req, handleReqFunc, deliverResFunc); err != nil {
		return err
	}
	return nil
}

func handleDeleteFile(s state.ConnectionState, m message.Message) error {
	req := m.(*message.DeleteFileRequest)
	deliverResFunc := func(res message.Response) error {
		writer := s.Conn
		headerBuffer := s.HeaderBuffer
		messageBuffer := s.Buffer
		return deliverSendRes(writer, headerBuffer, messageBuffer, res)
	}
	if err := handleReq(req, controller.HandleDeleteFileRequest, deliverResFunc); err != nil {
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

func deliverStreamRes(ctx context.Context, s state.ConnectionState, streamRes *service.StreamResponse) error {
	reader := streamRes.StreamReader
	defer closeReader(reader)

	var writer io.Writer = s.Conn
	headerBuffer := s.HeaderBuffer
	buffer := s.Buffer

	res := streamRes.StructResponse
	if err := deliverSendRes(writer, headerBuffer, buffer, res); err != nil {
		return err
	}

	toTransfer := streamRes.ToTransfer
	if err := transfer.Stream(ctx, reader, writer, buffer, toTransfer); err != nil {
		return err
	}
	return nil
}

func closeReader(reader io.Reader) {
	if closer, ok := reader.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			panic(err)
		}
	}
	panic("reader is not closer")
}

func deliverSendRes(writer io.Writer, headerBuffer []byte, messageBuffer *bytes.Buffer, res message.Response) error {
	m := res.(message.Message)
	if err := transfer.SendMessage(writer, headerBuffer, messageBuffer, m); err != nil {
		return err
	}
	return nil
}

const timeForRequest = 5 * time.Second

var ErrUnexpectedRequestType = errors.New("unexpected request type")
