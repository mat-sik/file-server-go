package client

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"log/slog"
	"net/http"
	"os"
)

func handelGetFileResponse(
	ctx context.Context,
	session netmsg.Session,
	filename string,
	res message.GetFileResponse,
) error {
	if res.Status != http.StatusOK {
		slog.Warn("GET file response:", "filename", filename, "status", res.Status)
		return nil
	}

	if err := downloadFile(ctx, session, filename, res.Size); err != nil {
		return err
	}

	slog.Info("GET file response:", "filename", filename, "status", res.Status, "size", res.Size)
	return nil
}

func downloadFile(
	ctx context.Context,
	session netmsg.Session,
	filename string,
	fileSize int,
) error {
	path := files.BuildClientFilePath(filename)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return session.StreamFromNet(ctx, file, fileSize)
}

func handlePutFileResponse(ctx context.Context, res message.PutFileResponse) {
	filename := filenameFromContextOrPanic(ctx)
	slog.Info("PUT file response:", "filename", filename, "status", res.Status)
}

func handleDeleteFileResponse(ctx context.Context, res message.DeleteFileResponse) {
	filename := filenameFromContextOrPanic(ctx)
	slog.Info("DELETE file response:", "filename", filename, "status", res.Status)
}

func handleGetFilenamesResponse(ctx context.Context, res message.GetFilenamesResponse) {
	pattern := patternFromContextOrPanic(ctx)
	if res.Status != http.StatusOK {
		slog.Warn("GET filenames response:", "pattern", pattern, "status", res.Status)
		return
	}
	slog.Info("GET filenames response:", "filenames", res.Filenames, "pattern", pattern, "status", res.Status)
}
