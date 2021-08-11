package message

import (
	"encoding/binary"
	"io"
)

type Writer struct {
	buff   []byte
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		buff:   make([]byte, 4),
		writer: w,
	}
}

// Write writes a pascal string. Not thread safe
func (w *Writer) Write(msg []byte) (int, error) {
	binary.BigEndian.PutUint32(w.buff, uint32(len(msg)))
	n, err := w.writer.Write(w.buff)
	if err != nil {
		return n, err
	}
	n, err = w.writer.Write(msg)
	return n, err
}
