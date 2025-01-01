package server

import (
	"context"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"github.com/mat-sik/file-server-go/internal/server/router"
	"net"
)

func RunServer(ctx context.Context, port int) error {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	defer safeListenerClose(listener)

	connCh := make(chan net.Conn)
	errCh := make(chan error)

	go acceptConnections(listener, connCh, errCh)

	for {
		select {
		case conn := <-connCh:
			go handleRequest(ctx, conn, errCh)
		case err = <-errCh:
			return err
		}
	}
}

func acceptConnections(listener net.Listener, connCh chan<- net.Conn, errCh chan<- error) {
	for {
		if conn, err := listener.Accept(); err != nil {
			errCh <- err
			return
		} else {
			connCh <- conn
		}
	}
}

func handleRequest(ctx context.Context, conn net.Conn, errCh chan<- error) {
	defer safeConnectionClose(conn)

	messageDispatcher := netmsg.NewSession(conn)
	serverRouter := router.ServerRouter{Session: messageDispatcher}

	if err := serverRouter.HandleRequest(ctx); err != nil {
		errCh <- err
	}
}

func safeConnectionClose(conn net.Conn) {
	if err := conn.Close(); err != nil {
		panic(err)
	}
}

func safeListenerClose(listener net.Listener) {
	if err := listener.Close(); err != nil {
		panic(err)
	}
}
