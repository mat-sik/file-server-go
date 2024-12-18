package decorated

import (
	"github.com/mat-sik/file-server-go/internal/message"
)

func New(res message.GetFileResponse, req *message.GetFileRequest) GetFileResponse {
	return GetFileResponse{
		GetFileResponse: res,
		Filename:        req.Filename,
	}
}

type GetFileResponse struct {
	message.GetFileResponse
	Filename string
}
