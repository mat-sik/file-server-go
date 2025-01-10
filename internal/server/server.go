package server

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"io"
	"log/slog"
	"net"
)

func Run(ctx context.Context, addr string) error {
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
			if !errors.Is(err, io.EOF) {
				return err
			}
			slog.Info("Connection closed from client")
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

	session := netmsg.NewSession(conn)
	handler := sessionHandler{Session: session}

	var err error
	for err == nil {
		err = handler.handleRequest(ctx)
	}

	errCh <- err
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
