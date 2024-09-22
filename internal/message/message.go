package message

import (
	"errors"
)

type Message interface {
	GetType() TypeName
}

type Request interface {
	GetRequestType() TypeName
}

type Response interface {
	GetResponseType() TypeName
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
	Filename string
}

func (req *GetFileRequest) GetType() TypeName {
	return GetFileResponseType
}

func (req *GetFileRequest) GetRequestType() TypeName {
	return req.GetType()
}

type GetFileResponse struct {
	Status int
	Size   int
}

func (res *GetFileResponse) GetType() TypeName {
	return GetFileRequestType
}

func (res *GetFileResponse) GetResponseType() TypeName {
	return res.GetType()
}

type PutFileRequest struct {
	FileName string
	Size     int
}

func (req *PutFileRequest) GetType() TypeName {
	return PutFileRequestType
}

func (req *PutFileRequest) GetRequestType() TypeName {
	return req.GetType()
}

type PutFileResponse struct {
	Status int
}

func (res *PutFileResponse) GetType() TypeName {
	return PutFileResponseType
}

func (res *PutFileResponse) GetResponseType() TypeName {
	return res.GetType()
}

type DeleteFileRequest struct {
	FileName string
}

func (req *DeleteFileRequest) GetType() TypeName {
	return DeleteFileRequestType
}

func (req *DeleteFileRequest) GetRequestType() TypeName {
	return req.GetType()
}

type DeleteFileResponse struct {
	Status int
}

func (res *DeleteFileResponse) GetType() TypeName {
	return DeleteFileResponseType
}

func (res *DeleteFileResponse) GetResponseType() TypeName {
	return res.GetType()
}
