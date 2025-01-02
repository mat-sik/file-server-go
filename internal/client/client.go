package client

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"net"
)

func RunClient(ctx context.Context, addr string, req message.Request) error {
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		return err
	}

	session := netmsg.NewSession(conn)
	sessionHandler := SessionHandler{Session: session}

	if err = sessionHandler.HandleRequest(ctx, req); err != nil {
		return err
	}

	return nil
}
