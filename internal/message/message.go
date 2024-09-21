package message

import (
	"errors"
)

type Request interface {
	GetRequestType() TypeName
}

type Response interface {
	GetResponseType() TypeName
}

type Holder struct {
	PayloadType   TypeName
	PayloadStruct any
}

type TypeName uint64

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

func (req *GetFileRequest) GetRequestType() TypeName {
	return GetFileRequestType
}

func NewGetFileRequestHolder(filename string) Holder {
	payload := GetFileRequest{Filename: filename}
	return Holder{
		PayloadType:   payload.GetRequestType(),
		PayloadStruct: payload,
	}
}

type GetFileResponse struct {
	Status int
	Size   int
}

func (res *GetFileResponse) GetResponseType() TypeName {
	return GetFileResponseType
}

func NewGetFileResponseHolder(status int, size int) Holder {
	payload := GetFileResponse{Status: status, Size: size}
	return Holder{
		PayloadType:   payload.GetResponseType(),
		PayloadStruct: payload,
	}
}

type PutFileRequest struct {
	FileName string
	Size     int
}

func (req *PutFileRequest) GetRequestType() TypeName {
	return PutFileRequestType
}

func NewPutFileRequestHolder(filename string, size int) Holder {
	payload := PutFileRequest{FileName: filename, Size: size}
	return Holder{
		PayloadType:   payload.GetRequestType(),
		PayloadStruct: payload,
	}
}

type PutFileResponse struct {
	Status int
}

func (res *PutFileResponse) GetResponseType() TypeName {
	return PutFileResponseType
}

func NewPutFileResponseHolder(status int) Holder {
	payload := PutFileResponse{Status: status}
	return Holder{
		PayloadType:   payload.GetResponseType(),
		PayloadStruct: payload,
	}
}

type DeleteFileRequest struct {
	FileName string
}

func (req *DeleteFileRequest) GetRequestType() TypeName {
	return DeleteFileRequestType
}

func NewDeleteFileRequestHolder(filename string) Holder {
	payload := DeleteFileRequest{FileName: filename}
	return Holder{
		PayloadType:   payload.GetRequestType(),
		PayloadStruct: payload,
	}
}

type DeleteFileResponse struct {
	Status int
}

func (res *DeleteFileResponse) GetResponseType() TypeName {
	return DeleteFileResponseType
}

func NewDeleteFileResponseHolder(status int) Holder {
	payload := DeleteFileResponse{Status: status}
	return Holder{
		PayloadType:   payload.GetResponseType(),
		PayloadStruct: payload,
	}
}
