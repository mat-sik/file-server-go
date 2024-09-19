package transfer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/message"
	"io"
	"testing"
)

func Test_sendMessage(t *testing.T) {
	//
	holder := message.NewPutFileRequestHolder("huge_file_name", 404)
	sizeBuffer := make([]byte, 12)
	messageBuffer := bytes.NewBuffer(make([]byte, 0, 1024))
	var socket io.Writer = bytes.NewBuffer(make([]byte, 0, 1024))

	err := sendMessage(socket, sizeBuffer, messageBuffer, &holder)
	if err != nil {
		t.Fatal(err)
	}

	messageBuffer.Reset()

	ctx := context.Background()
	out, err := receiveMessage(ctx, socket.(io.Reader), messageBuffer)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v", out)

	payload := out.PayloadStruct
	switch v := payload.(type) {
	case *message.PutFileRequest:
		fmt.Printf("%v", v)
	default:
		fmt.Printf("%v", v)
	}
}
