package response

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"log/slog"
	"net/http"
	"os"
)

func HandelGetFileResponse(
	ctx context.Context,
	session netmsg.Session,
	fileName string,
	res message.GetFileResponse,
) error {
	if res.Status != http.StatusOK {
		slog.Warn("GET file response:", "status", res.Status)
	}

	if err := handleGetFileResponse(ctx, session, fileName, res.Size); err != nil {
		return err
	}

	slog.Info("GET file response:", "status", res.Status, "size", res.Size)
	return nil
}

func handleGetFileResponse(
	ctx context.Context,
	session netmsg.Session,
	fileName string,
	fileSize int,
) error {
	path := files.BuildClientFilePath(fileName)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return session.StreamFromNet(ctx, file, fileSize)
}

func HandlePutFileResponse(res message.PutFileResponse) {
	slog.Info("PUT file response:", "status", res.Status)
}

func HandleDeleteFileResponse(res message.DeleteFileResponse) {
	slog.Info("DELETE file response:", "status", res.Status)
}
