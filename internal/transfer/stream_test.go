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
	buffer := limited.NewBuffer(make([]byte, 0, 4))

	reader := strings.NewReader("aaaabbbbcccc")
	writer := bytes.NewBuffer(make([]byte, 0, bytesBufferCap))

	expectedWriter := bytes.NewBuffer(make([]byte, 0, bytesBufferCap))
	expectedWriter.WriteString("aaaabbbb")
	// when
	ctx := context.Background()
	err := Stream(ctx, reader, writer, buffer, 8)
	// then
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(writer.Bytes(), expectedWriter.Bytes()) {
		t.Error(writer.Bytes(), expectedWriter.Bytes())
	}
}

const bytesBufferCap = 1024
