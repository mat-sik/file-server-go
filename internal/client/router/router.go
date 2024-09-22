package router

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/state"
	"io"
	"time"
)

func DeliverRequest(ctx context.Context, s state.ConnectionState, req message.Request) error {
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

const timeForRequest = 5 * time.Second
