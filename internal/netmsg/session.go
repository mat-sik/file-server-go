package netmsg

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"io"
	"net"
)

type Session struct {
	Conn   io.ReadWriteCloser
	Buffer []byte
}

func (s Session) SendMessage(msg message.Message) error {
	return sendMessage(msg, s.Buffer, s.Conn)
}

func (s Session) ReceiveMessage() (message.Message, error) {
	return receiveMessage(s.Conn, s.Buffer)
}

func (s Session) StreamToNet(ctx context.Context, reader io.Reader, toTransfer int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	limitedReader := io.LimitReader(reader, int64(toTransfer))
	_, err := io.CopyBuffer(s.Conn, limitedReader, s.Buffer)
	return err
}

func (s Session) StreamFromNet(ctx context.Context, writer io.Writer, toTransfer int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	limitedReader := io.LimitReader(s.Conn, int64(toTransfer))
	_, err := io.CopyBuffer(writer, limitedReader, s.Buffer)
	return err
}

func NewSession(conn net.Conn) Session {
	buffer := make([]byte, bufferSize)
	return Session{
		Conn:   conn,
		Buffer: buffer,
	}
}

const bufferSize = 4 * 1024
