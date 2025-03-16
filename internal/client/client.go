package client

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"net"
)

type Client struct {
	sessionHandler sessionHandler
}

func NewClient(addr string) (Client, error) {
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		return Client{}, err
	}

	session := netmsg.NewSession(conn)
	return Client{
		sessionHandler: sessionHandler{
			Session: session,
		},
	}, nil
}

func (c Client) Run(req message.Request) (message.Response, error) {
	ctx := context.Background()

	return c.sessionHandler.handleRequest(ctx, req)
}
