package decorated

import (
	"github.com/mat-sik/file-server-go/internal/message"
)

type GetFileResponse struct {
	*message.GetFileResponse
	FileName string
}
