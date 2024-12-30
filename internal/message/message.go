package message

import (
	"context"
	"errors"
	"github.com/mat-sik/file-server-go/internal/transfer/limited"
	"io"
)

type Message interface {
	GetType() TypeName
}

// StreamableMessage TODO: add reader that will represent the file either coming from tpc socket or os file system.
type StreamableMessage interface {
	Stream(ctx context.Context, writer io.Writer, headerBuffer []byte, messageBuffer *limited.Buffer) error
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
	Filename string
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

// PutFileRequest TODO: This type should have information about the amount of bytes of the file.
type PutFileRequest struct {
	Request
	FileName string
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
