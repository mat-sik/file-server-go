package enricher

import (
	"fmt"
	"github.com/mat-sik/file-server-go/internal/message"
)

func EnrichGetFileResponse(res message.Response, req message.Request) message.Response {
	getFileRequest, ok := req.(*message.GetFileRequest)
	if !ok {
		panic(fmt.Sprintf("GetFileRequest expected, received: %v", req))
	}
	fileName := getFileRequest.Filename

	getFileResponse := res.(*message.GetFileResponse)

	return &EnrichedGetFileResponse{
		Response: getFileResponse,
		Filename: fileName,
	}
}

type EnrichedGetFileResponse struct {
	message.Response
	Filename string
}
