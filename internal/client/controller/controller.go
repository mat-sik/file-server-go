package controller

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/client/service"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/state"
)

func GetFile(ctx context.Context, s state.ConnectionState, filename string) (message.Holder, error) {
	req := message.NewGetFileRequestHolder(filename)

	if err := service.SendRequest(s.Conn, s.HeaderBuffer, s.Buffer, &req); err != nil {
		return message.Holder{}, err
	}

	return service.HandleGetFile(ctx, s.Conn, s.Buffer, filename)
}

func PutFile(ctx context.Context, s state.ConnectionState, filename string) (message.Holder, error) {
	if err := service.PutFileHandleRequest(ctx, s.Conn, s.HeaderBuffer, s.Buffer, filename); err != nil {
		return message.Holder{}, err
	}

	return service.ReceiveResponse(ctx, s.Conn, s.Buffer)
}

func DeleteFile(ctx context.Context, s state.ConnectionState, filename string) (message.Holder, error) {
	req := message.NewDeleteFileRequestHolder(filename)

	if err := service.SendRequest(s.Conn, s.HeaderBuffer, s.Buffer, &req); err != nil {
		return message.Holder{}, err
	}

	return service.ReceiveResponse(ctx, s.Conn, s.Buffer)
}
