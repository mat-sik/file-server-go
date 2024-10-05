package conncontext

import (
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"github.com/mat-sik/file-server-go/internal/transfer/messheader"
	"net"
)

type ConnectionContext struct {
	Conn         net.Conn
	Buffer       *limited.Buffer
	HeaderBuffer []byte
}

func NewConnectionState(conn net.Conn) ConnectionContext {
	buffer := limited.NewBuffer(make([]byte, 0, bufferSize))
	headerBuffer := make([]byte, messheader.HeaderSize)
	return ConnectionContext{
		Conn:         conn,
		Buffer:       buffer,
		HeaderBuffer: headerBuffer,
	}
}

const bufferSize = 4 * 1024
