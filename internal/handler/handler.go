package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/controller"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"net"
)

type RequestState struct {
	Conn         net.Conn
	Buffer       *bytes.Buffer
	HeaderBuffer []byte
}

func HandleRequest(ctx context.Context, rs RequestState) error {
	holder, err := transfer.ReceiveMessage(ctx, rs.Conn, rs.Buffer)
	if err != nil {
		return err
	}
	switch holder.PayloadType {
	case message.GetFileRequestType:
		req := holder.PayloadStruct.(*message.GetFileRequest)
		if err = controller.GetFile(ctx, rs, *req); err != nil {
			return err
		}
	case message.PutFileRequestType:
		if err = handleRequest(ctx, rs, holder, controller.PutFile); err != nil {
			return err
		}
	case message.DeleteFileRequestType:
		if err = handleRequest(ctx, rs, holder, controller.DeleteFile); err != nil {
			return err
		}
	default:
		return ErrUnexpectedRequestType
	}
	return nil
}

func handleRequest[T any](
	ctx context.Context,
	rs RequestState,
	holder message.Holder,
	controllerFunc func(context.Context, RequestState, T) (message.Holder, error),
) error {
	req := holder.PayloadStruct.(*T)
	res, err := controllerFunc(ctx, rs, *req)
	if err != nil {
		return err
	}
	if err = transfer.SendMessage(rs.Conn, rs.HeaderBuffer, rs.Buffer, &res); err != nil {
		return err
	}
	return nil
}

var ErrUnexpectedRequestType = errors.New("unexpected request type")
