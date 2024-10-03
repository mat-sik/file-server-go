package transfer

import (
	"bytes"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/message"
	"io"
	"testing"
)

func Test_SendMessage_And_ReceiveMessage(t *testing.T) {
	//
	in := message.PutFileRequest{FileName: "huge_file_name", Size: 404}
	sizeBuffer := make([]byte, 12)
	messageBuffer := bytes.NewBuffer(make([]byte, 0, 1024))
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))

	var sendSocket io.Writer = buffer
	err := SendMessage(sendSocket, sizeBuffer, messageBuffer, &in)
	if err != nil {
		t.Fatal(err)
	}

	messageBuffer.Reset()

	var readSocket io.Reader = buffer
	out, err := ReceiveMessage(readSocket, messageBuffer)
	if err != nil {
		t.Fatal(err)
	}

	switch v := out.(type) {
	case *message.PutFileRequest:
		fmt.Printf("%v", v)
	default:
		fmt.Printf("%v", v)
	}
}
