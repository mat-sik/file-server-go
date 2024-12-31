package message

import (
	"errors"
)

type Message interface {
	GetType() TypeName
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

func (req GetFileRequest) GetType() TypeName {
	return GetFileRequestType
}

type GetFileResponse struct {
	Response
	Status int
	Size   int
}

func (res GetFileResponse) GetType() TypeName {
	return GetFileResponseType
}

type PutFileRequest struct {
	Request
	FileName string
	Size     int
}

func (req PutFileRequest) GetType() TypeName {
	return PutFileRequestType
}

type PutFileResponse struct {
	Response
	Status int
}

func (res PutFileResponse) GetType() TypeName {
	return PutFileResponseType
}

type DeleteFileRequest struct {
	Request
	FileName string
}

func (req DeleteFileRequest) GetType() TypeName {
	return DeleteFileRequestType
}

type DeleteFileResponse struct {
	Response
	Status int
}

func (res DeleteFileResponse) GetType() TypeName {
	return DeleteFileResponseType
}
