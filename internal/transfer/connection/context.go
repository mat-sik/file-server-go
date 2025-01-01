package connection

import (
	"github.com/mat-sik/file-server-go/internal/transfer/header"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
	"net"
)

type Context struct {
	io.ReadWriteCloser
	Buffer       *limited.Buffer
	HeaderBuffer []byte
}

func NewContext(conn net.Conn) Context {
	buffer := limited.NewBuffer(make([]byte, 0, bufferSize))
	headerBuffer := make([]byte, header.Size)
	return Context{
		ReadWriteCloser: conn,
		Buffer:          buffer,
		HeaderBuffer:    headerBuffer,
	}
}

const bufferSize = 4 * 1024
