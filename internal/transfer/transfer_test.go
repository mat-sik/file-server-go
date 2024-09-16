package transfer

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/message"
	"io"
	"strings"
	"testing"
)

func Test_transfer(t *testing.T) {
	// given
	buffer := bytes.NewBuffer(make([]byte, 0, 10))

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedBuffer := bytes.NewBuffer(make([]byte, 0, 1014))

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteString("one two three four f")

	// when
	ctx := context.Background()
	err := transfer(ctx, reader, writer, buffer, 20)
	// then
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(buffer.Bytes(), expectedBuffer.Bytes()) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_transfer_offset(t *testing.T) {
	// given
	buffer := bytes.NewBuffer(make([]byte, 2, 10))

	_, _ = buffer.ReadByte()
	_, _ = buffer.ReadByte()

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedBuffer := bytes.NewBuffer(make([]byte, 0, 1019))
	expectedBuffer.WriteString("our f")

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteString("one two three f")

	// when
	ctx := context.Background()
	err := transfer(ctx, reader, writer, buffer, 15)
	// then
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(buffer.Bytes(), expectedBuffer.Bytes()) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_transfer_buffered(t *testing.T) {
	// given
	buffer := bytes.NewBuffer(make([]byte, 2, 10))

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedBuffer := bytes.NewBuffer(make([]byte, 0, 1016))
	expectedBuffer.WriteString(" f")

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteByte(0)
	expectedWriter.WriteByte(0)
	expectedWriter.WriteString("one two three four")

	// when
	ctx := context.Background()
	err := transfer(ctx, reader, writer, buffer, 20)
	// then
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(buffer.Bytes(), expectedBuffer.Bytes()) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_transfer_offsetAndBufferedToTransferSmallerThanBuffer(t *testing.T) {
	// given
	buffer := bytes.NewBuffer(make([]byte, 2, 10))

	_, _ = buffer.ReadByte()

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedBuffer := bytes.NewBuffer(make([]byte, 0, 1021))
	expectedBuffer.WriteString("two th")

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteByte(0)
	expectedWriter.WriteString("one ")

	// when
	ctx := context.Background()
	err := transfer(ctx, reader, writer, buffer, 5)
	// then
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(buffer.Bytes(), expectedBuffer.Bytes()) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_transfer_offsetAndBuffered(t *testing.T) {
	// given
	buffer := bytes.NewBuffer(make([]byte, 2, 10))

	_, _ = buffer.ReadByte()

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedBuffer := bytes.NewBuffer(make([]byte, 0, 1015))
	expectedBuffer.WriteString("f")

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteByte(0)
	expectedWriter.WriteString("one two three four ")

	// when
	ctx := context.Background()
	err := transfer(ctx, reader, writer, buffer, 20)
	// then
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(buffer.Bytes(), expectedBuffer.Bytes()) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_sendMessage(t *testing.T) {
	//
	gob.Register(message.PutFileRequest{})

	holder := message.NewPutFileRequestHolder("huge_file_name", 404)
	sizeBuffer := make([]byte, 4)
	messageBuffer := bytes.NewBuffer(make([]byte, 0, 1024))
	var socket io.Writer = bytes.NewBuffer(make([]byte, 0, 1024))

	a := messageBuffer.Next(4)
	fmt.Printf("%v", a)

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
}

func Test_foo(t *testing.T) {
	holder := message.NewPutFileRequestHolder("huge_file_name", 404)

	out, err := json.Marshal(holder)
	if err != nil {
		t.Fatal(err)
	}

	var newHolder message.Holder
	if err = json.Unmarshal(out, &newHolder); err != nil {
		t.Fatal(err)
	}

	payload := newHolder.PayloadStruct.(map[string]any)
	var toDeserialize message.PutFileRequest
	err = message.UnmarshalType(payload, &toDeserialize)
	if err != nil {
		t.Fatal(err)
	}
}
