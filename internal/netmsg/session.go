package netmsg

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg/header"
	"github.com/mat-sik/file-server-go/internal/netmsg/limited"
	"io"
	"net"
)

type buffer interface {
	streamBuffer
	messageBuffer
}

type Session struct {
	Conn         io.ReadWriteCloser
	Buffer       buffer
	HeaderBuffer []byte
}

func (s Session) SendMessage(msg message.Message) error {
	return sendMessage(msg, s.HeaderBuffer, s.Buffer, s.Conn)
}

func (s Session) ReceiveMessage() (message.Message, error) {
	return receiveMessage(s.Conn, s.Buffer)
}

func (s Session) StreamToNet(ctx context.Context, reader io.Reader, toTransfer int) error {
	return stream(ctx, reader, s.Conn, s.Buffer, toTransfer)
}

func (s Session) StreamFromNet(ctx context.Context, writer io.Writer, toTransfer int) error {
	return stream(ctx, s.Conn, writer, s.Buffer, toTransfer)
}

func NewSession(conn net.Conn) Session {
	limitedBuffer := limited.NewBuffer(make([]byte, 0, bufferSize))
	headerBuffer := make([]byte, header.Size)
	return Session{
		Conn:         conn,
		Buffer:       limitedBuffer,
		HeaderBuffer: headerBuffer,
	}
}

const bufferSize = 4 * 1024
