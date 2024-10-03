package reshandler

import (
	"context"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/client/request/enricher"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/transfer/conncontext"
	"io"
)

func HandelGetFileResponse(ctx context.Context, connCtx conncontext.ConnectionContext, res message.Response) error {
	buffer := connCtx.Buffer
	defer buffer.Reset()

	enrichedGetFileResponse := res.(*enricher.EnrichedGetFileResponse)
	getFileResponse := enrichedGetFileResponse.Response.(*message.GetFileResponse)

	status := getFileResponse.Status
	if status != 200 {
		fmt.Printf("getFileResponse status: %d\n", status)
	}

	var reader io.Reader = connCtx.Conn
	fileName := enrichedGetFileResponse.Filename
	fileSize := getFileResponse.Size
	if err := handleGetFileResponse(ctx, reader, buffer, fileName, fileSize); err != nil {
		return err
	}
	fmt.Println(getFileResponse)

	return nil
}

func HandlePutFileResponse(res message.Response) {
	putFileResponse := res.(*message.PutFileResponse)

	status := putFileResponse.Status
	handlePutFileResponse(status)
}

func HandleDeleteFileResponse(res message.Response) {
	deleteFileResponse := res.(*message.DeleteFileResponse)

	status := deleteFileResponse.Status
	handleDeleteFileResponse(status)
}
