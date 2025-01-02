package message

import (
	"errors"
)

type Message interface {
	Type() TypeName
}

type Request interface {
	Message
}

type Response interface {
	Message
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

func TypeNameConverter(typeName TypeName) (Message, error) {
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
	Request
	FileName string
}

func (req GetFileRequest) Type() TypeName {
	return GetFileRequestType
}

type GetFileResponse struct {
	Response
	Status int
	Size   int
}

func NewGetFileResponse(status int, size int) *GetFileResponse {
	return &GetFileResponse{Status: status, Size: size}
}

func (res GetFileResponse) Type() TypeName {
	return GetFileResponseType
}

type PutFileRequest struct {
	Request
	FileName string
	Size     int
}

func (req PutFileRequest) Type() TypeName {
	return PutFileRequestType
}

type PutFileResponse struct {
	Response
	Status int
}

func NewPutFileResponse(status int) *PutFileResponse {
	return &PutFileResponse{Status: status}
}

func (res PutFileResponse) Type() TypeName {
	return PutFileResponseType
}

type DeleteFileRequest struct {
	Request
	FileName string
}

func (req DeleteFileRequest) Type() TypeName {
	return DeleteFileRequestType
}

type DeleteFileResponse struct {
	Response
	Status int
}

func NewDeleteFileResponse(status int) *DeleteFileResponse {
	return &DeleteFileResponse{Status: status}
}

func (res DeleteFileResponse) Type() TypeName {
	return DeleteFileResponseType
}
