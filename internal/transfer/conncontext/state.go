package conncontext

import (
	"bytes"
	"github.com/mat-sik/file-server-go/internal/transfer/messheader"
	"net"
)

type ConnectionContext struct {
	Conn         net.Conn
	Buffer       *bytes.Buffer
	HeaderBuffer []byte
}

func NewConnectionState(conn net.Conn) ConnectionContext {
	buffer := bytes.NewBuffer(make([]byte, 0, 4*1024))
	headerBuffer := make([]byte, messheader.HeaderSize)
	return ConnectionContext{
		Conn:         conn,
		Buffer:       buffer,
		HeaderBuffer: headerBuffer,
	}
}
