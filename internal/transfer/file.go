package transfer

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
	"os"
)

type FileStreamableMessage interface {
	GetMessage() message.Message
	GetFile() *os.File
	GetToTransfer() int
}

func StreamFromFile(
	ctx context.Context,
	writer io.Writer,
	headerBuffer []byte,
	messageBuffer *limited.Buffer,
	streamable FileStreamableMessage,
) error {
	defer messageBuffer.Reset()

	file := streamable.GetFile()
	defer safeFileClose(file)

	m := streamable.GetMessage()
	if err := SendMessage(writer, headerBuffer, messageBuffer, m); err != nil {
		return err
	}

	toTransfer := streamable.GetToTransfer()
	return Stream(ctx, file, writer, messageBuffer, toTransfer)
}

func safeFileClose(f *os.File) {
	if err := f.Close(); err != nil {
		panic(err)
	}
}
