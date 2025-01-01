package client

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/client/router"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"net"
)

func RunClient(ctx context.Context, hostname string) error {
	conn, err := net.Dial("tcp4", hostname)
	if err != nil {
		return err
	}

	messageDispatcher := netmsg.NewSession(conn)
	clientRouter := router.ClientRouter{Session: messageDispatcher}

	req := message.GetFileRequest{FileName: "foo.txt"}

	if err = clientRouter.HandleRequest(ctx, req); err != nil {
		return err
	}

	return nil
}
