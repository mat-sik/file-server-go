package message

import (
	"errors"
)

type Holder struct {
	PayloadType   TypeName
	PayloadStruct any
}

type TypeName int

const (
	GetFileRequestType TypeName = iota
	GetFileResponseType

	PutFileRequestType
	PutFileResponseType

	DeleteFileRequestType
	DeleteFileResponseType
)

func TypeNameConverter(typeName TypeName) (any, error) {
	switch typeName {
	case GetFileRequestType:
		return &GetFileRequest{}, nil
	case GetFileResponseType:
		return &GetFileResponse{}, nil
	case PutFileRequestType:
		return &PutFileRequest{}, nil
	case PutFileResponseType:
		return &PutFileResponse{}, nil
	case DeleteFileRequestType:
		return &DeleteFileRequest{}, nil
	case DeleteFileResponseType:
		return &DeleteFileResponse{}, nil
	default:
		return nil, errors.New("unknown type")
	}
}

type GetFileRequest struct {
	Filename string
}

func NewGetFileRequestHolder(filename string) Holder {
	return Holder{
		PayloadType:   GetFileRequestType,
		PayloadStruct: GetFileRequest{filename},
	}
}

var ErrFailedExtraction = errors.New("failed to extract payload")

type GetFileResponse struct {
	Status int
	Size   int
}

func NewGetFileResponseHolder(status int, size int) Holder {
	return Holder{
		PayloadType:   GetFileResponseType,
		PayloadStruct: GetFileResponse{Status: status, Size: size},
	}
}

type PutFileRequest struct {
	FileName string
	Size     int
}

func NewPutFileRequestHolder(filename string, size int) Holder {
	return Holder{
		PayloadType:   PutFileRequestType,
		PayloadStruct: PutFileRequest{filename, size},
	}
}

type PutFileResponse struct {
	Status int
}

func NewPutFileResponseHolder(status int) Holder {
	return Holder{
		PayloadType:   PutFileResponseType,
		PayloadStruct: PutFileResponse{status},
	}
}

type DeleteFileRequest struct {
	FileName string
}

func NewDeleteFileRequestHolder(filename string) Holder {
	return Holder{
		PayloadType:   DeleteFileRequestType,
		PayloadStruct: DeleteFileRequest{filename},
	}
}

type DeleteFileResponse struct {
	Status int
}

func NewDeleteFileResponseHolder(status int) Holder {
	return Holder{
		PayloadType:   DeleteFileResponseType,
		PayloadStruct: DeleteFileResponse{status},
	}
}
