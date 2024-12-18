package response

import (
	"context"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/client/request/decorated"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/connection"
	"io"
)

func HandelGetFileResponse(ctx context.Context, connCtx connection.Context, res decorated.GetFileResponse) error {
	buffer := connCtx.Buffer
	defer buffer.Reset()

	status := res.Status
	if status != 200 {
		fmt.Printf("getFileResponse status: %d\n", status)
	}

	var reader io.Reader = connCtx.Conn
	return handleGetFileResponse(ctx, reader, buffer, res.Filename, res.Size)
}

func HandlePutFileResponse(res message.PutFileResponse) {
	status := res.Status
	handlePutFileResponse(status)
}

func HandleDeleteFileResponse(res message.DeleteFileResponse) {
	status := res.Status
	handleDeleteFileResponse(status)
}
