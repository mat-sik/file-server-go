package main

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/client"
	"github.com/mat-sik/file-server-go/internal/server"
	"os"
)

func main() {
	args := os.Args
	opts := args[1]
	ctx := context.Background()
	if opts == "server" {
		if err := server.RunServer(ctx, 10000); err != nil {
			panic(err)
		}
	} else {
		if err := client.RunClient(ctx, "localhost:10000"); err != nil {
			panic(err)
		}
	}
}
