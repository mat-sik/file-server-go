package response

import (
	"context"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/envs"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/message/decorated"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"os"
	"path/filepath"
)

func HandelGetFileResponse(
	ctx context.Context,
	session netmsg.Session,
	res decorated.GetFileResponse,
) error {
	if res.Status != 200 {
		fmt.Printf("getFileResponse status: %d\n", res.Status)
	}
	if err := handleGetFileResponse(ctx, session, res.FileName, res.Size); err != nil {
		return err
	}

	fmt.Printf("handle get file response %d\n", res.Status)
	return nil
}

func handleGetFileResponse(
	ctx context.Context,
	session netmsg.Session,
	fileName string,
	fileSize int,
) error {
	defer session.Buffer.Reset()

	path := filepath.Join(envs.ClientDBPath, fileName)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return session.StreamFromNet(ctx, file, fileSize)
}

func HandlePutFileResponse(res message.PutFileResponse) {
	fmt.Printf("handle put file response %d\n", res.Status)
}

func HandleDeleteFileResponse(res message.DeleteFileResponse) {
	fmt.Printf("handle delete file response %d\n", res.Status)
}
