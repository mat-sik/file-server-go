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

type Session struct {
	Conn         io.ReadWriteCloser
	Buffer       Buffer
	HeaderBuffer []byte
}

func (d Session) SendMessage(m message.Message) error {
	return sendMessage(m, d.HeaderBuffer, d.Buffer, d.Conn)
}

func (d Session) ReceiveMessage() (message.Message, error) {
	return receiveMessage(d.Conn, d.Buffer)
}

func (d Session) StreamToNet(ctx context.Context, reader io.Reader, toTransfer int) error {
	return stream(ctx, reader, d.Conn, d.Buffer, toTransfer)
}

func (d Session) StreamFromNet(ctx context.Context, writer io.Writer, toTransfer int) error {
	return stream(ctx, d.Conn, writer, d.Buffer, toTransfer)
}

func NewSession(conn net.Conn) Session {
	buffer := limited.NewBuffer(make([]byte, 0, bufferSize))
	headerBuffer := make([]byte, header.Size)
	return Session{
		Conn:         conn,
		Buffer:       buffer,
		HeaderBuffer: headerBuffer,
	}
}

const bufferSize = 4 * 1024
