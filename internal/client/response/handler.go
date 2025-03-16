package response

import (
	"context"
	"github.com/mat-sik/file-server-go/internal/client/ctxvalue"
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
	filename string,
	res message.GetFileResponse,
) error {
	if res.Status != http.StatusOK {
		slog.Warn("GET file response:", "filename", filename, "status", res.Status)
		return nil
	}

	if err := handleGetFileResponse(ctx, session, filename, res.Size); err != nil {
		return err
	}

	slog.Info("GET file response:", "filename", filename, "status", res.Status, "size", res.Size)
	return nil
}

func handleGetFileResponse(
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

func HandlePutFileResponse(ctx context.Context, res message.PutFileResponse) {
	filename := filenameFromContext(ctx)
	slog.Info("PUT file response:", "filename", filename, "status", res.Status)
}

func HandleDeleteFileResponse(ctx context.Context, res message.DeleteFileResponse) {
	filename := filenameFromContext(ctx)
	slog.Info("DELETE file response:", "filename", filename, "status", res.Status)
}

func HandleGetFilenamesResponse(ctx context.Context, res message.GetFilenamesResponse) {
	pattern := patternFromContext(ctx)
	if res.Status != http.StatusOK {
		slog.Warn("GET filenames response:", "pattern", pattern, "status", res.Status)
		return
	}
	slog.Info("GET filenames response:", "filenames", res.Filenames, "pattern", pattern, "status", res.Status)
}

func patternFromContext(ctx context.Context) string {
	pattern, ok := ctxvalue.PatternFromContext(ctx)
	if !ok {
		panic("could not get pattern from context")
	}
	return pattern
}

func filenameFromContext(ctx context.Context) string {
	filename, ok := ctxvalue.FilenameFromContext(ctx)
	if !ok {
		panic("could not get filename from context")
	}
	return filename
}
