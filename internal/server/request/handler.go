package request

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"net/http"
	"os"
)

type Handler struct {
	files.Service
}

func NewHandler(fileService files.Service) Handler {
	return Handler{Service: fileService}
}

func (h Handler) HandleGetFileRequest(req message.GetFileRequest) (GetFileResponse, error) {
	fileHandle, ok := h.GetFile(req.FileName)
	if !ok {
		return GetFileResponse{
			GetFileResponse: message.GetFileResponse{
				Status: http.StatusNotFound,
				Size:   0,
			},
		}, nil
	}

	readLockedFile, err := fileHandle.NewReadLockedFile()
	if errors.Is(err, os.ErrNotExist) {
		return GetFileResponse{
			GetFileResponse: message.GetFileResponse{
				Status: http.StatusNotFound,
				Size:   0,
			},
		}, nil
	} else if err != nil {
		return GetFileResponse{}, err
	}

	fileSize, err := files.SizeOf(readLockedFile.File)
	if err != nil {
		defer files.LoggedClose(&readLockedFile)
		return GetFileResponse{}, err
	}

	return GetFileResponse{
		GetFileResponse: message.GetFileResponse{
			Status: http.StatusOK,
			Size:   fileSize,
		},
		ReadLockedFile: readLockedFile,
	}, nil
}

type GetFileResponse struct {
	message.GetFileResponse
	files.ReadLockedFile
}

func (h Handler) HandlePutFileRequest(
	ctx context.Context,
	session netmsg.Session,
	req message.PutFileRequest,
) (message.PutFileResponse, error) {
	saveFileFromNet := func(filename string) error {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer files.LoggedClose(file)

		return session.StreamFromNet(ctx, file, req.Size)
	}

	fileHandle := h.AddFile(req.FileName)
	if err := fileHandle.ExecuteWriteOP(saveFileFromNet); err != nil {
		return message.PutFileResponse{}, err
	}

	return message.PutFileResponse{
		Status: http.StatusCreated,
	}, nil
}

func (h Handler) HandleDeleteFileRequest(req message.DeleteFileRequest) (message.DeleteFileResponse, error) {
	err := h.RemoveFile(req.FileName)
	if errors.Is(err, os.ErrNotExist) {
		return message.DeleteFileResponse{
			Status: http.StatusNotFound,
		}, nil
	}
	if err != nil {
		return message.DeleteFileResponse{}, err
	}
	return message.DeleteFileResponse{
		Status: http.StatusOK,
	}, nil
}
