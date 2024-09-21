package router

import (
	"bytes"
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/server/controller"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/mheader"
	"net"
	"time"
)

type RequestState struct {
	Conn         net.Conn
	Buffer       *bytes.Buffer
	HeaderBuffer []byte
}

func NewRequestState(conn net.Conn) RequestState {
	buffer := bytes.NewBuffer(make([]byte, 4*1024))
	headerBuffer := make([]byte, mheader.HeaderSize)
	return RequestState{
		Conn:         conn,
		Buffer:       buffer,
		HeaderBuffer: headerBuffer,
	}
}

func RouteRequest(ctx context.Context, rs RequestState) error {
	ctx, cancel := context.WithTimeout(ctx, timeForRequest)
	defer cancel()

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
		putFileFunc := func(req message.PutFileRequest) (message.Holder, error) {
			return controller.PutFile(ctx, rs, req)
		}
		if err = handleReq(rs, holder, putFileFunc); err != nil {
			return err
		}
	case message.DeleteFileRequestType:
		if err = handleReq(rs, holder, controller.DeleteFile); err != nil {
			return err
		}
	default:
		return ErrUnexpectedRequestType
	}
	return nil
}

func handleReq[T any](
	rs RequestState,
	holder message.Holder,
	reqFunc func(T) (message.Holder, error),
) error {
	req := holder.PayloadStruct.(*T)
	res, err := reqFunc(*req)
	if err != nil {
		return err
	}
	if err = transfer.SendMessage(rs.Conn, rs.HeaderBuffer, rs.Buffer, &res); err != nil {
		return err
	}
	return nil
}

const timeForRequest = 5 * time.Second

var ErrUnexpectedRequestType = errors.New("unexpected request type")
