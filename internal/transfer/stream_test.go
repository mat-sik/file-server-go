package transfer

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"strings"
	"testing"
	"time"
)

func Test_Stream(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		buffer        StreamableBuffer
		reader        io.Reader
		writer        *bytes.Buffer
		mockFunc      func(StreamableBuffer, context.Context, io.Reader, io.Writer)
		assertFunc    func(StreamableBuffer, context.Context, io.Reader, io.Writer)
		toTransfer    int
		wantedData    string
		expectedError error
	}{
		{
			name:   "normal case",
			ctx:    context.Background(),
			buffer: &MockStreamableBuffer{},
			reader: strings.NewReader("aaaabbbbcccc"),
			writer: bytes.NewBuffer(make([]byte, 0, bytesBufferCap)),
			mockFunc: func(b StreamableBuffer, _ context.Context, r io.Reader, w io.Writer) {
				m, _ := b.(*MockStreamableBuffer)

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
			assertFunc: func(b StreamableBuffer, _ context.Context, _ io.Reader, _ io.Writer) {
				m, _ := b.(*MockStreamableBuffer)
				m.AssertExpectations(t)
			},
			toTransfer: 8,
			wantedData: "aaaabbbb",
		},
		{
			name:   "exact buffer size",
			ctx:    context.Background(),
			buffer: &MockStreamableBuffer{},
			reader: strings.NewReader("aaaabbbb"),
			writer: bytes.NewBuffer(make([]byte, 0, bytesBufferCap)),
			mockFunc: func(b StreamableBuffer, _ context.Context, r io.Reader, w io.Writer) {
				m, _ := b.(*MockStreamableBuffer)

				len0 := m.On("Len").Return(0).Once()
				reset0 := m.On("Reset").Return().Once().NotBefore(len0)
				read0 := m.On("SingleReadFrom", r).Return(4, nil).Once().NotBefore(reset0)
				len1 := m.On("Len").Return(4).Once().NotBefore(read0)
				m.On("SingleWriteTo", w, 4).Return(4, nil, []byte("aaaa")).Once().NotBefore(len1)
			},
			assertFunc: func(b StreamableBuffer, _ context.Context, _ io.Reader, _ io.Writer) {
				m, _ := b.(*MockStreamableBuffer)
				m.AssertExpectations(t)
			},
			toTransfer: 4,
			wantedData: "aaaa",
		},
		{
			name:   "zero transfer",
			ctx:    context.Background(),
			buffer: &MockStreamableBuffer{},
			reader: strings.NewReader("aaaabbbb"),
			writer: bytes.NewBuffer(make([]byte, 0, bytesBufferCap)),
			mockFunc: func(_ StreamableBuffer, _ context.Context, _ io.Reader, _ io.Writer) {
			},
			assertFunc: func(b StreamableBuffer, _ context.Context, _ io.Reader, _ io.Writer) {
				m, _ := b.(*MockStreamableBuffer)
				m.AssertExpectations(t)
			},
			toTransfer: 0,
			wantedData: "",
		},
		{
			name:   "ctx early return",
			ctx:    &MockContext{},
			buffer: &MockStreamableBuffer{},
			reader: strings.NewReader("aaaabbbb"),
			writer: bytes.NewBuffer(make([]byte, 0, bytesBufferCap)),
			mockFunc: func(_ StreamableBuffer, c context.Context, _ io.Reader, _ io.Writer) {
				m, _ := c.(*MockContext)

				doneCh := make(chan struct{}, 1)
				doneCh <- struct{}{}

				done0 := m.On("Done").Return(doneCh).Once()
				m.On("Err").Return(context.Canceled).Once().NotBefore(done0)
			},
			assertFunc: func(buffer StreamableBuffer, ctx context.Context, _ io.Reader, _ io.Writer) {
				mockBuffer, _ := buffer.(*MockStreamableBuffer)
				mockCtx, _ := ctx.(*MockContext)

				mockCtx.AssertExpectations(t)
				mockBuffer.AssertExpectations(t)
			},
			toTransfer:    4,
			wantedData:    "",
			expectedError: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			want := []byte(tt.wantedData)

			tt.mockFunc(tt.buffer, tt.ctx, tt.reader, tt.writer)

			// when
			err := Stream(tt.ctx, tt.reader, tt.writer, tt.buffer, tt.toTransfer)
			// then
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)

				got := tt.writer.Bytes()
				assert.Equal(t, want, got)
			}

			tt.assertFunc(tt.buffer, tt.ctx, tt.reader, tt.writer)
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

type MockContext struct {
	mock.Mock
}

func (m *MockContext) Deadline() (deadline time.Time, ok bool) {
	args := m.Called()
	return args.Get(0).(time.Time), args.Bool(1)
}

func (m *MockContext) Done() <-chan struct{} {
	args := m.Called()
	return args.Get(0).(chan struct{})
}

func (m *MockContext) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockContext) Value(key interface{}) interface{} {
	args := m.Called(key)
	return args.Get(0)
}

const bytesBufferCap = 1024
