package netmsg

import (
	"fmt"
	"github.com/mat-sik/file-server-go/internal/generated/netmsgpb"
	"github.com/mat-sik/file-server-go/internal/message"
	"google.golang.org/protobuf/proto"
	"io"
)

func sendMessage(msg message.Message, buffer []byte, writer io.Writer) error {
	wrapperMsg := toProto(msg)

	msgBytes, err := proto.Marshal(&wrapperMsg)
	if err != nil {
		return err
	}

	msgSize := uint32(len(msgBytes))
	msgHeader := header{
		payloadSize: msgSize,
	}
	if err = encodeHeader(msgHeader, buffer); err != nil {
		return err
	}

	if _, err = writer.Write(buffer[:headerSize]); err != nil {
		return err
	}
	if _, err = writer.Write(msgBytes); err != nil {
		return err
	}
	return nil
}

func receiveMessage(reader io.Reader, buffer []byte) (message.Message, error) {
	limitedReader := io.LimitReader(reader, int64(headerSize))
	if _, err := limitedReader.Read(buffer); err != nil {
		return nil, err
	}

	msgHeader, err := decodeHeader(buffer)
	if err != nil {
		return nil, err
	}

	limitedReader = io.LimitReader(reader, int64(msgHeader.payloadSize))
	if _, err = limitedReader.Read(buffer); err != nil {
		return nil, err
	}

	msg := &netmsgpb.MessageWrapper{}
	if err = proto.Unmarshal(buffer[:msgHeader.payloadSize], msg); err != nil {
		return nil, err
	}

	return fromProto(msg), nil
}

func toProto(msg message.Message) netmsgpb.MessageWrapper {
	switch msg := msg.(type) {
	case message.GetFileRequest:
		return netmsgpb.MessageWrapper{
			Message: &netmsgpb.MessageWrapper_GetFileRequest{
				GetFileRequest: &netmsgpb.GetFileRequest{
					Filename: &msg.Filename,
				},
			},
		}
	case message.GetFileResponse:
		status := int32(msg.Status)
		size := int64(msg.Size)
		return netmsgpb.MessageWrapper{
			Message: &netmsgpb.MessageWrapper_GetFileResponse{
				GetFileResponse: &netmsgpb.GetFileResponse{
					Status: &status,
					Size:   &size,
				},
			},
		}
	case message.PutFileRequest:
		size := int64(msg.Size)
		return netmsgpb.MessageWrapper{
			Message: &netmsgpb.MessageWrapper_PutFileRequest{
				PutFileRequest: &netmsgpb.PutFileRequest{
					Filename: &msg.Filename,
					Size:     &size,
				},
			},
		}
	case message.PutFileResponse:
		status := int32(msg.Status)
		return netmsgpb.MessageWrapper{
			Message: &netmsgpb.MessageWrapper_PutFileResponse{
				PutFileResponse: &netmsgpb.PutFileResponse{
					Status: &status,
				},
			},
		}
	case message.DeleteFileRequest:
		return netmsgpb.MessageWrapper{
			Message: &netmsgpb.MessageWrapper_DeleteFileRequest{
				DeleteFileRequest: &netmsgpb.DeleteFileRequest{
					Filename: &msg.Filename,
				},
			},
		}
	case message.DeleteFileResponse:
		status := int32(msg.Status)
		return netmsgpb.MessageWrapper{
			Message: &netmsgpb.MessageWrapper_DeleteFileResponse{
				DeleteFileResponse: &netmsgpb.DeleteFileResponse{
					Status: &status,
				},
			},
		}
	case message.GetFilenamesRequest:
		return netmsgpb.MessageWrapper{
			Message: &netmsgpb.MessageWrapper_GetFilenamesRequest{
				GetFilenamesRequest: &netmsgpb.GetFilenamesRequest{
					MatchRegex: &msg.MatchRegex,
				},
			},
		}
	case message.GetFilenamesResponse:
		status := int32(msg.Status)
		return netmsgpb.MessageWrapper{
			Message: &netmsgpb.MessageWrapper_GetFilenamesResponse{
				GetFilenamesResponse: &netmsgpb.GetFilenamesResponse{
					Status:   &status,
					Filename: msg.Filenames,
				},
			},
		}
	default:
		panic(fmt.Sprintf("unexpected message type %T", msg))
	}
}

func fromProto(wrapper *netmsgpb.MessageWrapper) message.Message {
	switch msg := wrapper.GetMessage().(type) {
	case *netmsgpb.MessageWrapper_GetFileRequest:
		req := msg.GetFileRequest
		return message.GetFileRequest{
			Filename: req.GetFilename(),
		}
	case *netmsgpb.MessageWrapper_GetFileResponse:
		req := msg.GetFileResponse
		return message.GetFileResponse{
			Status: int(req.GetStatus()),
			Size:   int(req.GetSize()),
		}
	case *netmsgpb.MessageWrapper_PutFileRequest:
		req := msg.PutFileRequest
		return message.PutFileRequest{
			Filename: req.GetFilename(),
			Size:     int(req.GetSize()),
		}
	case *netmsgpb.MessageWrapper_PutFileResponse:
		req := msg.PutFileResponse
		return message.PutFileResponse{
			Status: int(req.GetStatus()),
		}
	case *netmsgpb.MessageWrapper_DeleteFileRequest:
		req := msg.DeleteFileRequest
		return message.DeleteFileRequest{
			Filename: req.GetFilename(),
		}
	case *netmsgpb.MessageWrapper_DeleteFileResponse:
		req := msg.DeleteFileResponse
		return message.DeleteFileResponse{
			Status: int(req.GetStatus()),
		}
	case *netmsgpb.MessageWrapper_GetFilenamesRequest:
		req := msg.GetFilenamesRequest
		return message.GetFilenamesRequest{
			MatchRegex: req.GetMatchRegex(),
		}
	case *netmsgpb.MessageWrapper_GetFilenamesResponse:
		req := msg.GetFilenamesResponse
		return message.GetFilenamesResponse{
			Status:    int(req.GetStatus()),
			Filenames: req.GetFilename(),
		}
	default:
		panic(fmt.Sprintf("unexpected message type %T", msg))
	}
}
