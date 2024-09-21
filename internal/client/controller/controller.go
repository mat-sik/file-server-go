package controller

import (
	"bytes"
	"context"
	"github.com/mat-sik/file-server-go/internal/client/service"
	"github.com/mat-sik/file-server-go/internal/message"
	"net"
)

type RequestState struct {
	Conn          net.Conn
	HeaderBuffer  []byte
	MessageBuffer *bytes.Buffer
}

func GetFile(ctx context.Context, rs RequestState, filename string) (message.Holder, error) {
	req := message.NewGetFileRequestHolder(filename)

	if err := service.SendRequest(rs.Conn, rs.HeaderBuffer, rs.MessageBuffer, &req); err != nil {
		return message.Holder{}, err
	}

	return service.HandleGetFile(ctx, rs.Conn, rs.MessageBuffer, filename)
}

func PutFile(ctx context.Context, rs RequestState, filename string) (message.Holder, error) {
	if err := service.PutFileHandleRequest(ctx, rs.Conn, rs.HeaderBuffer, rs.MessageBuffer, filename); err != nil {
		return message.Holder{}, err
	}

	return service.ReceiveResponse(ctx, rs.Conn, rs.MessageBuffer)
}

func DeleteFile(ctx context.Context, rs RequestState, filename string) (message.Holder, error) {
	req := message.NewDeleteFileRequestHolder(filename)

	if err := service.SendRequest(rs.Conn, rs.HeaderBuffer, rs.MessageBuffer, &req); err != nil {
		return message.Holder{}, err
	}

	return service.ReceiveResponse(ctx, rs.Conn, rs.MessageBuffer)
}
