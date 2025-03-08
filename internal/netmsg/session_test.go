package netmsg

import (
	"bytes"
	"github.com/mat-sik/file-server-go/internal/message"
	"reflect"
	"testing"
)

func Test_should_SendMessage_And_ReceiveIt(t *testing.T) {
	testCases := []struct {
		name    string
		message message.Message
	}{
		{name: "PUT File Request", message: message.PutFileRequest{FileName: "foo.txt", Size: 404}},
		{name: "PUT File Response", message: message.PutFileResponse{Status: 200}},
		{name: "GET File Request", message: message.GetFileRequest{FileName: "foo.txt"}},
		{name: "GET File Response", message: message.GetFileResponse{Status: 200, Size: 404}},
		{name: "DELETE File Request", message: message.DeleteFileRequest{FileName: "foo.txt"}},
		{name: "DELETE File Response", message: message.DeleteFileResponse{Status: 200}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer := make([]byte, 1024)
			readWriteCloser := &mockReadWriteCloser{Buffer: *bytes.NewBuffer(make([]byte, 0, 1024))}

			session := Session{
				Conn:   readWriteCloser,
				Buffer: buffer,
			}

			if err := session.SendMessage(tc.message); err != nil {
				t.Fatal(err)
			}

			out, err := session.ReceiveMessage()
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(out, tc.message) {
				t.Fatalf("got %T of value %v, want %T of value %v", out, out, tc.message, tc.message)
			}
		})
	}
}

type mockReadWriteCloser struct {
	bytes.Buffer
}

func (mock *mockReadWriteCloser) Close() error {
	return nil
}
