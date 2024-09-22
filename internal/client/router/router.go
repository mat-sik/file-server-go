package router

import (
	"bytes"
	"context"
	"github.com/mat-sik/file-server-go/internal/client/service"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/state"
	"io"
)

func deliverStreamReq(ctx context.Context, s state.ConnectionState, streamReq *service.StreamRequest) error {
	reader := streamReq.Reader
	defer closeReader(reader)

	var writer io.Writer = s.Conn
	headerBuffer := s.HeaderBuffer
	buffer := s.Buffer

	req := streamReq.StructRequest
	if err := deliverReq(writer, headerBuffer, buffer, req); err != nil {
		return err
	}

	toTransfer := streamReq.ToTransfer
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

func deliverReq(writer io.Writer, headerBuffer []byte, messageBuffer *bytes.Buffer, req message.Request) error {
	m := req.(message.Message)
	if err := transfer.SendMessage(writer, headerBuffer, messageBuffer, m); err != nil {
		return err
	}
	return nil
}
