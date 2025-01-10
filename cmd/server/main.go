package main

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/server"
)

//go:generate protoc --proto_path=./../.. --go_out=./../.. --go_opt=module=github.com/mat-sik/file-server-go netmsg.proto
func main() {
	ctx := context.Background()
	if err := server.Run(ctx, ":44696"); err != nil {
		panic(err)
	}
}
