package router

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/server/controller"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/state"
	"time"
)

func RouteRequest(ctx context.Context, s state.ConnectionState) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

	holder, err := transfer.ReceiveMessage(ctx, s.Conn, s.Buffer)
	if err != nil {
		return err
	}
	switch holder.PayloadType {
	case message.GetFileRequestType:
		req := holder.PayloadStruct.(*message.GetFileRequest)
		if err = controller.GetFile(ctx, s, *req); err != nil {
			return err
		}
	case message.PutFileRequestType:
		putFileFunc := func(req message.PutFileRequest) (message.Holder, error) {
			return controller.PutFile(ctx, s, req)
		}
		if err = handleReq(s, holder, putFileFunc); err != nil {
			return err
		}
	case message.DeleteFileRequestType:
		if err = handleReq(s, holder, controller.DeleteFile); err != nil {
			return err
		}
	default:
		return ErrUnexpectedRequestType
	}
	return nil
}

func handleReq[T any](
	s state.ConnectionState,
	holder message.Holder,
	reqFunc func(T) (message.Holder, error),
) error {
	req := holder.PayloadStruct.(*T)
	res, err := reqFunc(*req)
	if err != nil {
		return err
	}
	if err = transfer.SendMessage(s.Conn, s.HeaderBuffer, s.Buffer, &res); err != nil {
		return err
	}
	return nil
}

const timeForRequest = 5 * time.Second

var ErrUnexpectedRequestType = errors.New("unexpected request type")
