package client

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/client/router"
	"github.com/mat-sik/file-server-go/internal/client/service"
	"github.com/mat-sik/file-server-go/internal/transfer/state"
	"net"
)

func RunClient(ctx context.Context, hostname string) error {
	conn, err := net.Dial("tcp", hostname)
	if err != nil {
		return err
	}

	s := state.NewConnectionState(conn)

	req, err := service.HandleGetFileRequest("foo.txt")
	if err != nil {
		return err
	}

	if err = router.HandleRequest(ctx, s, req); err != nil {
		return err
	}

	return nil
}
