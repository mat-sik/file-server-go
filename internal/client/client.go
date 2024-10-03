package client

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/client/request"
	"github.com/mat-sik/file-server-go/internal/client/router"
	"github.com/mat-sik/file-server-go/internal/transfer/conncontext"
	"net"
)

func RunClient(ctx context.Context, hostname string) error {
	conn, err := net.Dial("tcp4", hostname)
	if err != nil {
		return err
	}

	connCtx := conncontext.NewConnectionState(conn)

	req, err := request.NewGetFileRequest("foo.txt")
	if err != nil {
		return err
	}

	if err = router.HandleRequest(ctx, connCtx, req); err != nil {
		return err
	}

	return nil
}
