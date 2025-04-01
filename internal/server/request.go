package server

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/message"
	"github.com/mat-sik/file-server-go/internal/netmsg"
	"net/http"
	"os"
	"regexp"
)

type handler struct {
	files.Service
}

func newHandler(fileService files.Service) handler {
	return handler{Service: fileService}
}

func (h handler) handleGetFileRequest(req message.GetFileRequest) (getFileResponse, error) {
	fileHandle, ok := h.GetFile(req.Filename)
	if !ok {
		return getFileResponse{
			GetFileResponse: message.GetFileResponse{
				Status: http.StatusNotFound,
				Size:   0,
			},
		}, nil
	}

	readLockedFile, err := fileHandle.NewReadLockedFile()
	if errors.Is(err, os.ErrNotExist) {
		return getFileResponse{
			GetFileResponse: message.GetFileResponse{
				Status: http.StatusNotFound,
				Size:   0,
			},
		}, nil
	} else if err != nil {
		return getFileResponse{}, err
	}

	fileSize, err := files.SizeOf(readLockedFile.File)
	if err != nil {
		defer files.LoggedClose(&readLockedFile)
		return getFileResponse{}, err
	}

	return getFileResponse{
		GetFileResponse: message.GetFileResponse{
			Status: http.StatusOK,
			Size:   fileSize,
		},
		ReadLockedFile: readLockedFile,
	}, nil
}

type getFileResponse struct {
	message.GetFileResponse
	files.ReadLockedFile
}

func (h handler) handlePutFileRequest(
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

	fileHandle := h.AddFile(req.Filename)
	if err := fileHandle.ExecuteWriteOP(saveFileFromNet); err != nil {
		return message.PutFileResponse{}, err
	}

	return message.PutFileResponse{
		Status: http.StatusCreated,
	}, nil
}

func (h handler) handleDeleteFileRequest(req message.DeleteFileRequest) (message.DeleteFileResponse, error) {
	err := h.RemoveFile(req.Filename)
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

func (h handler) handleGetFilenamesRequest(req message.GetFilenamesRequest) (message.GetFilenamesResponse, error) {
	pattern, err := regexp.Compile(req.MatchRegex)
	if err != nil {
		return message.GetFilenamesResponse{
			Status: http.StatusBadRequest,
		}, nil
	}

	var filteredFilenames []string
	for _, filename := range h.GetAllFilenames() {
		if pattern.MatchString(filename) {
			filteredFilenames = append(filteredFilenames, filename)
		}
	}

	return message.GetFilenamesResponse{
		Status:    http.StatusOK,
		Filenames: filteredFilenames,
	}, nil
}
