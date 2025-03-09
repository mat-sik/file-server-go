package message

type GetFileRequest struct {
	Filename string
}

type GetFileResponse struct {
	Status int
	Size   int
}

type PutFileRequest struct {
	Filename string
	Size     int
}

type PutFileResponse struct {
	Status int
}

type DeleteFileRequest struct {
	Filename string
}

type DeleteFileResponse struct {
	Status int
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
