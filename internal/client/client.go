package client

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"net"
)

func RunClient(ctx context.Context, hostname string) error {
	conn, err := net.Dial("tcp4", hostname)
	if err != nil {
		return err
	}

	session := netmsg.NewSession(conn)
	sessionHandler := SessionHandler{Session: session}

	req := message.GetFileRequest{FileName: "foo.txt"}

	if err = sessionHandler.HandleRequest(ctx, req); err != nil {
		return err
	}

	return nil
}
