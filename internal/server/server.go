package server

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"io"
	"log/slog"
	"net"
	"sync"
)

func Run(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	defer files.LoggedClose(listener)

	return run(ctx, listener)
}

func RunWithWaitGroup(ctx context.Context, wg *sync.WaitGroup, addr string) error {
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	defer files.LoggedClose(listener)
	wg.Done()

	return run(ctx, listener)
}

func run(ctx context.Context, listener net.Listener) error {
	connCh := make(chan net.Conn)
	errCh := make(chan error)

	go acceptConnections(listener, connCh, errCh)

	return connectionLoop(ctx, connCh, errCh)
}

func acceptConnections(listener net.Listener, connCh chan<- net.Conn, errCh chan<- error) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			errCh <- err
			return
		}
		connCh <- conn
	}
}

func connectionLoop(ctx context.Context, connCh <-chan net.Conn, errCh chan error) error {
	for {
		select {
		case conn := <-connCh:
			go handleRequest(ctx, conn, errCh)
		case err := <-errCh:
			if !errors.Is(err, io.EOF) {
				return err
			}
			slog.Info("Connection closed from client")
		case <-ctx.Done():
			return nil
		}
	}
}

func handleRequest(ctx context.Context, conn net.Conn, errCh chan<- error) {
	defer files.LoggedClose(conn)

	session := netmsg.NewSession(conn)

	fileService := files.NewService()
	requestHandler := newHandler(fileService)

	sh := sessionHandler{
		session: session,
		handler: requestHandler,
	}

	var err error
	for err == nil {
		err = sh.handleRequest(ctx)
	}

	errCh <- err
}
