package main

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/server"
)

func main() {
	ctx := context.Background()
	if err := server.RunServer(ctx, ":44696"); err != nil {
		panic(err)
	}
}
