package netmsg

import (
	"bytes"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg/limited"
	"testing"
)

func Test_SendMessage_And_ReceiveMessage(t *testing.T) {
	//
	in := message.PutFileRequest{FileName: "huge_file_name", Size: 404}
	sizeBuffer := make([]byte, 12)
	messageBuffer := limited.NewBuffer(make([]byte, 0, 1024))
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))

	readWriteCloser := &mockReadWriteCloser{Buffer: *buffer}

	messageDispatcher := MessageDispatcher{
		Conn:         readWriteCloser,
		Buffer:       messageBuffer,
		HeaderBuffer: sizeBuffer,
	}

	err := messageDispatcher.SendMessage(in)
	if err != nil {
		t.Fatal(err)
	}

	messageBuffer.Reset()

	out, err := messageDispatcher.ReceiveMessage()
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

type mockReadWriteCloser struct {
	bytes.Buffer
}

func (mock *mockReadWriteCloser) Close() error {
	return nil
}
