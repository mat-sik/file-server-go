package client

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"net"
)

type Client struct {
	sessionHandler SessionHandler
}

func NewClient(addr string) (Client, error) {
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		return Client{}, err
	}

	session := netmsg.NewSession(conn)
	sessionHandler := SessionHandler{Session: session}
	return Client{sessionHandler: sessionHandler}, nil
}

func (c Client) Run(req message.Request) error {
	ctx := context.Background()

	return c.sessionHandler.HandleRequest(ctx, req)
}
