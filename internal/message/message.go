package message

type GetFileRequest struct {
	FileName string
}

type GetFileResponse struct {
	Status int
	Size   int
}

func NewGetFileResponse(status int, size int) *GetFileResponse {
	return &GetFileResponse{Status: status, Size: size}
}

type PutFileRequest struct {
	FileName string
	Size     int
}

type PutFileResponse struct {
	Status int
}

func NewPutFileResponse(status int) *PutFileResponse {
	return &PutFileResponse{Status: status}
}

type DeleteFileRequest struct {
	FileName string
}

type DeleteFileResponse struct {
	Status int
}

func NewDeleteFileResponse(status int) *DeleteFileResponse {
	return &DeleteFileResponse{Status: status}
}

type Message interface {
	isMessage()
}

func (_ GetFileRequest) isMessage() {
}

func (_ GetFileResponse) isMessage() {
}

func (_ PutFileRequest) isMessage() {
}

func (_ PutFileResponse) isMessage() {
}

func (_ DeleteFileRequest) isMessage() {
}

func (_ DeleteFileResponse) isMessage() {
}

type Request interface {
	isMessage()
	isRequest()
}

func (_ GetFileRequest) isRequest() {
}

func (_ PutFileRequest) isRequest() {
}

func (_ DeleteFileRequest) isRequest() {
}

type Response interface {
	isMessage()
	isResponse()
}

func (_ GetFileResponse) isResponse() {
}

func (_ PutFileResponse) isResponse() {
}

func (_ DeleteFileResponse) isResponse() {
}
