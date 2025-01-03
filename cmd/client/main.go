package main

import (
	"github.com/mat-sik/file-server-go/internal/client"
	"github.com/mat-sik/file-server-go/internal/message"
)

func main() {
	webClient, err := client.NewClient(":44696")
	if err != nil {
		panic(err)
	}

	req := &message.GetFileRequest{FileName: "foo.txt"}

	err = webClient.Run(req)
	if err != nil {
		panic(err)
	}
}
