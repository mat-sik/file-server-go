package transfer

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func Test_stream(t *testing.T) {
	// given
	buffer := bytes.NewBuffer(make([]byte, 0, 10))

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedBuffer := bytes.NewBuffer(make([]byte, 0, 1014))

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteString("one two three four f")

	// when
	ctx := context.Background()
	err := stream(ctx, reader, writer, buffer, 20)
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

func Test_stream_offset(t *testing.T) {
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
	err := stream(ctx, reader, writer, buffer, 15)
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

func Test_stream_buffered(t *testing.T) {
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
	err := stream(ctx, reader, writer, buffer, 20)
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

func Test_stream_offsetAndBufferedToTransferSmallerThanBuffer(t *testing.T) {
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
	err := stream(ctx, reader, writer, buffer, 5)
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

func Test_stream_offsetAndBuffered(t *testing.T) {
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
	err := stream(ctx, reader, writer, buffer, 20)
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
