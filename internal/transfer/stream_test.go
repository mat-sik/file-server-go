package transfer

import (
	"bytes"
	"context"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"strings"
	"testing"
)

func Test_Stream(t *testing.T) {
	// given
	buffer := limited.NewBuffer(make([]byte, 0, 10))

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteString("one two three four f")

	// when
	ctx := context.Background()
	err := Stream(ctx, reader, writer, buffer, 20)
	// then
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_Stream_offset(t *testing.T) {
	// given
	buffer := limited.NewBuffer(make([]byte, 2, 10))
	setOffset(buffer, 2)

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteString("one two three f")

	// when
	ctx := context.Background()
	err := Stream(ctx, reader, writer, buffer, 15)
	// then
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_Stream_buffered(t *testing.T) {
	// given
	buffer := limited.NewBuffer(make([]byte, 2, 10))

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteByte(0)
	expectedWriter.WriteByte(0)
	expectedWriter.WriteString("one two three four")

	// when
	ctx := context.Background()
	err := Stream(ctx, reader, writer, buffer, 20)
	// then
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_Stream_offsetAndBufferedToTransferSmallerThanBuffer(t *testing.T) {
	// given
	buffer := limited.NewBuffer(make([]byte, 2, 10))
	setOffset(buffer, 1)

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteByte(0)
	expectedWriter.WriteString("one ")

	// when
	ctx := context.Background()
	err := Stream(ctx, reader, writer, buffer, 5)
	// then
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_Stream_offsetAndBuffered(t *testing.T) {
	// given
	buffer := limited.NewBuffer(make([]byte, 2, 10))
	setOffset(buffer, 1)

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteByte(0)
	expectedWriter.WriteString("one two three four ")

	// when
	ctx := context.Background()
	err := Stream(ctx, reader, writer, buffer, 20)
	// then
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func setOffset(buffer *limited.Buffer, n int) {
	_, _ = buffer.Read(make([]byte, n))
}
