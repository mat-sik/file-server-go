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
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))

	var sendSocket io.Writer = buffer
	err := sendMessage(sendSocket, sizeBuffer, messageBuffer, &holder)
	if err != nil {
		t.Fatal(err)
	}

	messageBuffer.Reset()

	var readSocket io.Reader = buffer
	ctx := context.Background()
	out, err := receiveMessage(ctx, readSocket, messageBuffer)
	if err != nil {
		t.Fatal(err)
	}

	payload := out.PayloadStruct
	switch v := payload.(type) {
	case *message.PutFileRequest:
		fmt.Printf("%v", v)
	default:
		fmt.Printf("%v", v)
	}
}
