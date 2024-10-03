package reshandler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/envs"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"io"
	"os"
	"path/filepath"
)

func handleGetFileResponse(
	ctx context.Context,
	reader io.Reader,
	buffer *bytes.Buffer,
	fileName string,
	fileSize int,
) error {
	path := filepath.Join(envs.ClientDBPath, fileName)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	var writer io.Writer = file
	return transfer.Stream(ctx, reader, writer, buffer, fileSize)
}

func handlePutFileResponse(status int) {
	fmt.Printf("handle put file response %d\n", status)
}

func handleDeleteFileResponse(status int) {
	fmt.Printf("handle delete file response %d\n", status)
}
