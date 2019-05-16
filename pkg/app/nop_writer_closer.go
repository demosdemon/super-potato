package app

import "io"

type NopWriterCloser struct {
	io.Writer
}

func (NopWriterCloser) Close() error {
	return nil
}

func NewNopWriterCloser(w io.Writer) (io.WriteCloser, error) {
	return NopWriterCloser{w}, nil
}
