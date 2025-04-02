package netmsg

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"io"
	"net"
)

type Session struct {
	conn   io.ReadWriteCloser
	buffer []byte
}

func (s Session) SendMessage(msg message.Message) error {
	return sendMessage(msg, s.buffer, s.conn)
}

func (s Session) ReceiveMessage() (message.Message, error) {
	return receiveMessage(s.conn, s.buffer)
}

func (s Session) StreamToNet(ctx context.Context, reader io.Reader, toTransfer int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	limitedReader := io.LimitReader(reader, int64(toTransfer))
	_, err := io.CopyBuffer(s.conn, limitedReader, s.buffer)
	return err
}

func (s Session) StreamFromNet(ctx context.Context, writer io.Writer, toTransfer int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	limitedReader := io.LimitReader(s.conn, int64(toTransfer))
	_, err := io.CopyBuffer(writer, limitedReader, s.buffer)
	return err
}

func NewSession(conn net.Conn) Session {
	buffer := make([]byte, bufferSize)
	return Session{
		conn:   conn,
		buffer: buffer,
	}
}

const bufferSize = 4 * 1024
