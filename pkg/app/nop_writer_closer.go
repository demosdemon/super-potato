package app

import "io"

type nopWriterCloser struct {
	io.Writer
}

func (nopWriterCloser) Close() error {
	return nil
}

func NopWriterCloser(w io.Writer) (io.WriteCloser, error) {
	return nopWriterCloser{w}, nil
}
