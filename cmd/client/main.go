package main

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/client"
	"github.com/mat-sik/file-server-go/internal/message"
)

func main() {
	ctx := context.Background()
	req := message.PutFileRequest{FileName: "foo.txt", Size: 4}
	if err := client.RunClient(ctx, ":44696", req); err != nil {
		panic(err)
	}
}
