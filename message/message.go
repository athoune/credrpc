package message

import (
	"io"
)

type ReadWriter struct {
	r *Reader
	w *Writer
}

func NewReaderWriter(rw io.ReadWriter) *ReadWriter {
	return &ReadWriter{
		r: NewReader(rw),
		w: NewWriter(rw),
	}
}

func (rw *ReadWriter) Write(msg []byte) (int, error) {
	return rw.w.Write(msg)
}

func (rw *ReadWriter) Read() ([]byte, error) {
	return rw.Read()
}
