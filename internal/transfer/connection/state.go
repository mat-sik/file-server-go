package connection

import (
	"github.com/mat-sik/file-server-go/internal/transfer/header"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"net"
)

type Context struct {
	Conn         net.Conn
	Buffer       *limited.Buffer
	HeaderBuffer []byte
}

func NewContext(conn net.Conn) Context {
	buffer := limited.NewBuffer(make([]byte, 0, bufferSize))
	headerBuffer := make([]byte, header.Size)
	return Context{
		Conn:         conn,
		Buffer:       buffer,
		HeaderBuffer: headerBuffer,
	}
}

const bufferSize = 4 * 1024
