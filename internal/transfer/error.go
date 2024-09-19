package transfer

import "errors"

var (
	ErrTooBigMessage        = errors.New("buffer is too small to fit the message")
	ErrHeaderBufferTooSmall = errors.New("mheader buffer too small")
)
