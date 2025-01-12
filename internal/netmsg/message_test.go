package netmsg

import (
	"bytes"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/message"
	"testing"
)

func Test_SendMessage_And_ReceiveMessage(t *testing.T) {
	in := message.PutFileRequest{FileName: "huge_file_name", Size: 404}
	buffer := make([]byte, 1024)
	bytesBuffer := bytes.NewBuffer(make([]byte, 0, 1024))

	readWriteCloser := &mockReadWriteCloser{Buffer: *bytesBuffer}

	messageDispatcher := Session{
		Conn:   readWriteCloser,
		Buffer: buffer,
	}

	err := messageDispatcher.SendMessage(in)
	if err != nil {
		t.Fatal(err)
	}

	out, err := messageDispatcher.ReceiveMessage()
	if err != nil {
		t.Fatal(err)
	}

	switch v := out.(type) {
	case message.PutFileRequest:
		fmt.Printf("%v", v)
	default:
		t.Errorf("%v", v)
	}
}

type mockReadWriteCloser struct {
	bytes.Buffer
}

func (mock *mockReadWriteCloser) Close() error {
	return nil
}
