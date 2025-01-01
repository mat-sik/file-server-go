package response

import (
	"context"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/envs"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/message/decorated"
	"github.com/mat-sik/file-server-go/internal/transfer"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
	"os"
	"path/filepath"
)

func HandelGetFileResponse(
	ctx context.Context,
	reader io.Reader,
	buffer *limited.Buffer,
	res decorated.GetFileResponse,
) error {
	if res.Status != 200 {
		fmt.Printf("getFileResponse status: %d\n", res.Status)
	}
	return handleGetFileResponse(ctx, reader, buffer, res.FileName, res.Size)
}

func handleGetFileResponse(
	ctx context.Context,
	reader io.Reader,
	buffer *limited.Buffer,
	fileName string,
	fileSize int,
) error {
	defer buffer.Reset()

	path := filepath.Join(envs.ClientDBPath, fileName)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	return transfer.Stream(ctx, reader, file, buffer, fileSize)
}

func HandlePutFileResponse(res message.PutFileResponse) {
	fmt.Printf("handle put file response %d\n", res.Status)
}

func HandleDeleteFileResponse(res message.DeleteFileResponse) {
	fmt.Printf("handle delete file response %d\n", res.Status)
}
