package message

import (
	"bytes"
	"context"
	"errors"
	"io"
)

type Message interface {
	GetType() TypeName
}

type StreamableMessage interface {
	Stream(ctx context.Context, writer io.Writer, headerBuffer []byte, messageBuffer *bytes.Buffer) error
}

type Request interface {
	GetRequestType() RequestTypeName
}

type Response interface {
	GetResponseType() ResponseTypeName
}

type TypeName uint64

const (
	GetFileRequestTypeNum TypeName = iota
	GetFileResponseTypeNum

	PutFileRequestTypeNum
	PutFileResponseTypeNum

	DeleteFileRequestTypeNum
	DeleteFileResponseTypeNum
)

type RequestTypeName TypeName

const (
	GetFileRequestType    = RequestTypeName(GetFileRequestTypeNum)
	PutFileRequestType    = RequestTypeName(PutFileRequestTypeNum)
	DeleteFileRequestType = RequestTypeName(DeleteFileRequestTypeNum)
)

type ResponseTypeName TypeName

const (
	GetFileResponseType    = ResponseTypeName(GetFileResponseTypeNum)
	PutFileResponseType    = ResponseTypeName(PutFileResponseTypeNum)
	DeleteFileResponseType = ResponseTypeName(DeleteFileResponseTypeNum)
)

func TypeNameConverter(typeName TypeName) (Message, error) {
	switch typeName {
	case GetFileRequestTypeNum:
		return &GetFileRequest{}, nil
	case GetFileResponseTypeNum:
		return &GetFileResponse{}, nil
	case PutFileRequestTypeNum:
		return &PutFileRequest{}, nil
	case PutFileResponseTypeNum:
		return &PutFileResponse{}, nil
	case DeleteFileRequestTypeNum:
		return &DeleteFileRequest{}, nil
	case DeleteFileResponseTypeNum:
		return &DeleteFileResponse{}, nil
	default:
		return nil, errors.New("unknown type")
	}
}

type GetFileRequest struct {
	Filename string
}

func (req *GetFileRequest) GetType() TypeName {
	return GetFileRequestTypeNum
}

func (req *GetFileRequest) GetRequestType() RequestTypeName {
	return RequestTypeName(req.GetType())
}

type GetFileResponse struct {
	Status int
	Size   int
}

func (res *GetFileResponse) GetType() TypeName {
	return GetFileResponseTypeNum
}

func (res *GetFileResponse) GetResponseType() ResponseTypeName {
	return ResponseTypeName(res.GetType())
}

type PutFileRequest struct {
	FileName string
	Size     int
}

func (req *PutFileRequest) GetType() TypeName {
	return PutFileRequestTypeNum
}

func (req *PutFileRequest) GetRequestType() RequestTypeName {
	return RequestTypeName(req.GetType())
}

type PutFileResponse struct {
	Status int
}

func (res *PutFileResponse) GetType() TypeName {
	return PutFileResponseTypeNum
}

func (res *PutFileResponse) GetResponseType() ResponseTypeName {
	return ResponseTypeName(res.GetType())
}

type DeleteFileRequest struct {
	FileName string
}

func (req *DeleteFileRequest) GetType() TypeName {
	return DeleteFileRequestTypeNum
}

func (req *DeleteFileRequest) GetRequestType() RequestTypeName {
	return RequestTypeName(req.GetType())
}

type DeleteFileResponse struct {
	Status int
}

func (res *DeleteFileResponse) GetType() TypeName {
	return DeleteFileResponseTypeNum
}

func (res *DeleteFileResponse) GetResponseType() ResponseTypeName {
	return ResponseTypeName(res.GetType())
}
