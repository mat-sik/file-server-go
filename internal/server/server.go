package server

import (
	"context"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/server/router"
	"github.com/mat-sik/file-server-go/internal/transfer/conncontext"
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

	connCtx := conncontext.NewConnectionState(conn)

	if err := router.HandleRequest(ctx, connCtx); err != nil {
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
