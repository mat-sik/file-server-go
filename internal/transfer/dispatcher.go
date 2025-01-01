package transfer

import (
	"github.com/mat-sik/file-server-go/internal/transfer/header"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
	"net"
)

type MessageDispatcher struct {
	Conn         io.ReadWriteCloser
	Buffer       Buffer
	HeaderBuffer []byte
}

type Buffer interface {
	Streamer
	Messenger
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
