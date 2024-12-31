package client

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/client/router"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/connection"
	"net"
)

func RunClient(ctx context.Context, hostname string) error {
	conn, err := net.Dial("tcp4", hostname)
	if err != nil {
		return err
	}

	connCtx := connection.NewContext(conn)

	req := message.GetFileRequest{FileName: "foo.txt"}

	if err = router.HandleRequest(ctx, connCtx, req); err != nil {
		return err
	}

	return nil
}
