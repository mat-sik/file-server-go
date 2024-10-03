package transfer

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func Test_Stream(t *testing.T) {
	// given
	buffer := bytes.NewBuffer(make([]byte, 0, 10))

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedBuffer := bytes.NewBuffer(make([]byte, 0, 1014))

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteString("one two three four f")

	// when
	ctx := context.Background()
	err := Stream(ctx, reader, writer, buffer, 20)
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

func Test_Stream_offset(t *testing.T) {
	// given
	buffer := bytes.NewBuffer(make([]byte, 2, 10))
	setOffset(buffer, 2)

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedBuffer := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedBuffer.WriteString("ree f")
	setOffset(expectedBuffer, 5)

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteString("one two three f")

	// when
	ctx := context.Background()
	err := Stream(ctx, reader, writer, buffer, 15)
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

func Test_Stream_buffered(t *testing.T) {
	// given
	buffer := bytes.NewBuffer(make([]byte, 2, 10))

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedBuffer := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedBuffer.WriteString("three four")
	setOffset(expectedBuffer, 10)

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
	if !bytes.Equal(buffer.Bytes(), expectedBuffer.Bytes()) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_Stream_offsetAndBufferedToTransferSmallerThanBuffer(t *testing.T) {
	// given
	buffer := bytes.NewBuffer(make([]byte, 2, 10))
	setOffset(buffer, 1)

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedBuffer := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedBuffer.WriteString("one ")
	setOffset(expectedBuffer, 4)

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
	if !bytes.Equal(buffer.Bytes(), expectedBuffer.Bytes()) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_Stream_offsetAndBuffered(t *testing.T) {
	// given
	buffer := bytes.NewBuffer(make([]byte, 2, 10))
	setOffset(buffer, 1)

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedBuffer := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedBuffer.WriteString("ree four ")
	setOffset(expectedBuffer, 9)

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
	if !bytes.Equal(buffer.Bytes(), expectedBuffer.Bytes()) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func setOffset(buffer *bytes.Buffer, n int) {
	for range n {
		_, _ = buffer.ReadByte()
	}
}
