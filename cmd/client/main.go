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

	getFileReq := &message.GetFileRequest{FileName: "foo.txt"}

	err = webClient.Run(getFileReq)
	if err != nil {
		panic(err)
	}

	delFileReq := &message.DeleteFileRequest{FileName: "foo.txt"}

	err = webClient.Run(delFileReq)
	if err != nil {
		panic(err)
	}

	putFileReq := &message.PutFileRequest{FileName: "foo.txt"}

	err = webClient.Run(putFileReq)
	if err != nil {
		panic(err)
	}
}
