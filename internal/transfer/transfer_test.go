package transfer

import (
	"bytes"
	"strings"
	"testing"
)

func Test_transfer(t *testing.T) {
	// given
	buffer := make([]byte, 10)

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedOffset := 10
	expectedBuffered := 0

	expectedBuffer := make([]byte, 10)
	copy(expectedBuffer, "ree four f")

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteString("one two three four f")

	// when
	offset, buffered, err := transfer(reader, writer, buffer, 0, 0, 20)
	// then
	if err != nil {
		t.Error(err)
	}
	if offset != expectedOffset {
		t.Error(offset, expectedOffset)
	}
	if buffered != expectedBuffered {
		t.Error(buffered, expectedBuffered)
	}
	if !bytes.Equal(buffer, expectedBuffer) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_transfer_offset(t *testing.T) {
	// given
	buffer := make([]byte, 10)
	buffer[0] = 1
	buffer[1] = 2

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedOffset := 5
	expectedBuffered := 5

	expectedBuffer := make([]byte, 10)
	copy(expectedBuffer, "ree four f")

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteString("one two three f")

	// when
	offset, buffered, err := transfer(reader, writer, buffer, 2, 0, 15)
	// then
	if err != nil {
		t.Error(err)
	}
	if offset != expectedOffset {
		t.Error(offset, expectedOffset)
	}
	if buffered != expectedBuffered {
		t.Error(buffered, expectedBuffered)
	}
	if !bytes.Equal(buffer, expectedBuffer) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_transfer_buffered(t *testing.T) {
	// given
	buffer := make([]byte, 10)
	buffer[0] = 1
	buffer[1] = 2

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedOffset := 8
	expectedBuffered := 2

	expectedBuffer := make([]byte, 10)
	copy(expectedBuffer, "ree four f")

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteByte(1)
	expectedWriter.WriteByte(2)
	expectedWriter.WriteString("one two three four")

	// when
	offset, buffered, err := transfer(reader, writer, buffer, 0, 2, 20)
	// then
	if err != nil {
		t.Error(err)
	}
	if offset != expectedOffset {
		t.Error(offset, expectedOffset)
	}
	if buffered != expectedBuffered {
		t.Error(buffered, expectedBuffered)
	}
	if !bytes.Equal(buffer, expectedBuffer) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_transfer_offsetAndBufferedToTransferSmallerThanBuffer(t *testing.T) {
	// given
	buffer := make([]byte, 10)
	buffer[0] = 1
	buffer[1] = 2

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedOffset := 4
	expectedBuffered := 6

	expectedBuffer := make([]byte, 10)
	copy(expectedBuffer, "one two th")

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteByte(2)
	expectedWriter.WriteString("one ")

	// when
	offset, buffered, err := transfer(reader, writer, buffer, 1, 1, 5)
	// then
	if err != nil {
		t.Error(err)
	}
	if offset != expectedOffset {
		t.Error(offset, expectedOffset)
	}
	if buffered != expectedBuffered {
		t.Error(buffered, expectedBuffered)
	}
	if !bytes.Equal(buffer, expectedBuffer) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

func Test_transfer_offsetAndBuffered(t *testing.T) {
	// given
	buffer := make([]byte, 10)
	buffer[0] = 1
	buffer[1] = 2

	reader := strings.NewReader("one two three four five six")
	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	expectedOffset := 9
	expectedBuffered := 1

	expectedBuffer := make([]byte, 10)
	copy(expectedBuffer, "ree four f")

	expectedWriter := bytes.NewBuffer(make([]byte, 0, 1024))
	expectedWriter.WriteByte(2)
	expectedWriter.WriteString("one two three four ")

	// when
	offset, buffered, err := transfer(reader, writer, buffer, 1, 1, 20)
	// then
	if err != nil {
		t.Error(err)
	}
	if offset != expectedOffset {
		t.Error(offset, expectedOffset)
	}
	if buffered != expectedBuffered {
		t.Error(buffered, expectedBuffered)
	}
	if !bytes.Equal(buffer, expectedBuffer) {
		t.Error(buffer, expectedBuffer)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}
