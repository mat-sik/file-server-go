package message

import (
	"errors"
	"fmt"
	"reflect"
)

type Holder struct {
	TypeName      TypeName
	PayloadStruct any
}

type TypeName string

const (
	GetFileRequestType  TypeName = "GetFileReq"
	GetFileResponseType TypeName = "GetFileRes"

	PutFileRequestType  TypeName = "PutFileReq"
	PutFileResponseType TypeName = "PutFileRes"

	DeleteFileRequestType  TypeName = "DeleteFileReq"
	DeleteFileResponseType TypeName = "DeleteFileRes"
)

func ExtractType[T any](holder Holder) (T, error) {
	payload, ok := holder.PayloadStruct.(T)
	if !ok {
		return payload, ErrFailedExtraction
	}
	return payload, nil
}

func UnmarshalType(serializedStruct map[string]any, passedStruct any) error {
	value := reflect.ValueOf(passedStruct).Elem()
	valueType := value.Type()

	for i := range valueType.NumField() {
		field := valueType.Field(i)
		fieldValue := value.Field(i)

		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set %s field value", field.Name)
		}

		mapValue, exists := serializedStruct[field.Name]
		if !exists {
			return fmt.Errorf("cannot find %s field value", field.Name)
		}

		if mapValueReflect := reflect.ValueOf(mapValue); mapValueReflect.Type().ConvertibleTo(fieldValue.Type()) {
			fieldValue.Set(mapValueReflect.Convert(fieldValue.Type()))
		} else {
			return fmt.Errorf("cannot assign value of type %T to field %s", mapValue, field.Name)
		}
	}
	return nil
}

type GetFileRequest struct {
	Filename string
}

func NewGetFileRequestHolder(filename string) Holder {
	return Holder{
		TypeName:      GetFileRequestType,
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
		TypeName:      GetFileResponseType,
		PayloadStruct: GetFileResponse{Status: status, Size: size},
	}
}

type PutFileRequest struct {
	FileName string
	Size     int
}

func NewPutFileRequestHolder(filename string, size int) Holder {
	return Holder{
		TypeName:      PutFileRequestType,
		PayloadStruct: PutFileRequest{filename, size},
	}
}

type PutFileResponse struct {
	Status int
}

func NewPutFileResponseHolder(status int) Holder {
	return Holder{
		TypeName:      PutFileResponseType,
		PayloadStruct: PutFileResponse{status},
	}
}

type DeleteFileRequest struct {
	FileName string
}

func NewDeleteFileRequestHolder(filename string) Holder {
	return Holder{
		TypeName:      DeleteFileRequestType,
		PayloadStruct: DeleteFileRequest{filename},
	}
}

type DeleteFileResponse struct {
	Status int
}

func NewDeleteFileResponseHolder(status int) Holder {
	return Holder{
		TypeName:      DeleteFileResponseType,
		PayloadStruct: DeleteFileResponse{status},
	}
}
