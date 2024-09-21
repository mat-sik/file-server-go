package state

import (
	"bytes"
	"github.com/mat-sik/file-server-go/internal/transfer/mheader"
	"net"
)

type ConnectionState struct {
	Conn         net.Conn
	Buffer       *bytes.Buffer
	HeaderBuffer []byte
}

func NewConnectionState(conn net.Conn) ConnectionState {
	buffer := bytes.NewBuffer(make([]byte, 4*1024))
	headerBuffer := make([]byte, mheader.HeaderSize)
	return ConnectionState{
		Conn:         conn,
		Buffer:       buffer,
		HeaderBuffer: headerBuffer,
	}
}
