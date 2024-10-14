package transfer

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"strings"
	"testing"
)

func Test_Stream(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(buffer *MockStreamableBuffer, r io.Reader, w io.Writer)
		readerData  string
		toTransfer  int
		wantedData  string
		expectError bool
	}{
		{
			name: "normal case",
			mockFunc: func(m *MockStreamableBuffer, r io.Reader, w io.Writer) {
				len0 := m.On("Len").Return(0).Once()
				reset0 := m.On("Reset").Return().Once().NotBefore(len0)
				read0 := m.On("SingleReadFrom", r).Return(4, nil).Once().NotBefore(reset0)
				len1 := m.On("Len").Return(4).Once().NotBefore(read0)
				write0 := m.On("SingleWriteTo", w, 4).Return(4, nil, []byte("aaaa")).Once().NotBefore(len1)
				reset1 := m.On("Reset").Return().Once().NotBefore(write0)
				read1 := m.On("SingleReadFrom", r).Return(4, nil).Once().NotBefore(reset1)
				len2 := m.On("Len").Return(4).Once().NotBefore(read1)
				m.On("SingleWriteTo", w, 4).Return(4, nil, []byte("bbbb")).Once().NotBefore(len2)
			},
			readerData: "aaaabbbbcccc",
			toTransfer: 8,
			wantedData: "aaaabbbb",
		},
		{
			name: "exact buffer size",
			mockFunc: func(m *MockStreamableBuffer, r io.Reader, w io.Writer) {
				len0 := m.On("Len").Return(0).Once()
				reset0 := m.On("Reset").Return().Once().NotBefore(len0)
				read0 := m.On("SingleReadFrom", r).Return(4, nil).Once().NotBefore(reset0)
				len1 := m.On("Len").Return(4).Once().NotBefore(read0)
				m.On("SingleWriteTo", w, 4).Return(4, nil, []byte("aaaa")).Once().NotBefore(len1)
			},
			readerData: "aaaabbbb",
			toTransfer: 4,
			wantedData: "aaaa",
		},
		{
			name: "zero transfer",
			mockFunc: func(m *MockStreamableBuffer, r io.Reader, w io.Writer) {
			},
			readerData: "aaaabbbb",
			toTransfer: 0,
			wantedData: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			reader := strings.NewReader(tt.readerData)
			writer := bytes.NewBuffer(make([]byte, 0, bytesBufferCap))

			want := []byte(tt.wantedData)
			ctx := context.Background()

			mockBuffer := &MockStreamableBuffer{}
			tt.mockFunc(mockBuffer, reader, writer)

			// when
			err := Stream(ctx, reader, writer, mockBuffer, tt.toTransfer)
			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				got := writer.Bytes()
				assert.Equal(t, want, got)
			}

			mockBuffer.AssertExpectations(t)
		})
	}
}

type MockStreamableBuffer struct {
	mock.Mock
}

func (m *MockStreamableBuffer) SingleReadFrom(reader io.Reader) (int, error) {
	args := m.Called(reader)
	n := args.Int(0)
	_, _ = reader.Read(make([]byte, n))
	return n, args.Error(1)
}

func (m *MockStreamableBuffer) SingleWriteTo(writer io.Writer, limit int) (int, error) {
	args := m.Called(writer, limit)
	_, _ = writer.Write(args.Get(2).([]byte))
	return args.Int(0), args.Error(1)
}

func (m *MockStreamableBuffer) Len() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockStreamableBuffer) Reset() {
	_ = m.Called()
}

const bytesBufferCap = 1024
