package main

import (
	"github.com/mat-sik/file-server-go/internal/client"
	"github.com/mat-sik/file-server-go/internal/message"
)

//go:generate protoc --proto_path=./../.. --go_out=./../.. --go_opt=module=github.com/mat-sik/file-server-go netmsg.proto
func main() {
	webClient, err := client.NewClient(":44696")
	if err != nil {
		panic(err)
	}

	getFileReq := message.GetFileRequest{Filename: "foo.txt"}

	_, err = webClient.Run(getFileReq)
	if err != nil {
		panic(err)
	}

	delFileReq := message.DeleteFileRequest{Filename: "foo.txt"}

	_, err = webClient.Run(delFileReq)
	if err != nil {
		panic(err)
	}

	putFileReq := message.PutFileRequest{Filename: "foo.txt"}

	_, err = webClient.Run(putFileReq)
	if err != nil {
		panic(err)
	}
}
