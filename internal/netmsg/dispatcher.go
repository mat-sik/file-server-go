package netmsg

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg/header"
	"github.com/mat-sik/file-server-go/internal/netmsg/limited"
	"io"
	"net"
)

type Buffer interface {
	StreamBuffer
	MessageBuffer
}

type MessageDispatcher struct {
	Conn         io.ReadWriteCloser
	Buffer       Buffer
	HeaderBuffer []byte
}

func (d MessageDispatcher) SendMessage(m message.Message) error {
	return sendMessage(m, d.HeaderBuffer, d.Buffer, d.Conn)
}

func (d MessageDispatcher) ReceiveMessage() (message.Message, error) {
	return receiveMessage(d.Conn, d.Buffer)
}

func (d MessageDispatcher) StreamToNet(ctx context.Context, reader io.Reader, toTransfer int) error {
	return stream(ctx, reader, d.Conn, d.Buffer, toTransfer)
}

func (d MessageDispatcher) StreamFromNet(ctx context.Context, writer io.Writer, toTransfer int) error {
	return stream(ctx, d.Conn, writer, d.Buffer, toTransfer)
}

func NewMessageDispatcher(conn net.Conn) MessageDispatcher {
	buffer := limited.NewBuffer(make([]byte, 0, bufferSize))
	headerBuffer := make([]byte, header.Size)
	return MessageDispatcher{
		Conn:         conn,
		Buffer:       buffer,
		HeaderBuffer: headerBuffer,
	}
}

const bufferSize = 4 * 1024
