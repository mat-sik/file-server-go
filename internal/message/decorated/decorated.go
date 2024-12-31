package decorated

import (
	"github.com/mat-sik/file-server-go/internal/message"
)

func New(res message.GetFileResponse, fileName string) GetFileResponse {
	return GetFileResponse{
		GetFileResponse: res,
		FileName:        fileName,
	}
}

type GetFileResponse struct {
	message.GetFileResponse
	FileName string
}
