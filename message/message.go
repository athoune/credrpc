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

type Reader struct {
	buff   []byte
	reader io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		buff:   make([]byte, 4),
		reader: r,
	}
}

func (r *Reader) Read() ([]byte, error) {
	_, err := io.ReadFull(r.reader, r.buff)
	if err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(r.buff)
	resp := make([]byte, size)
	_, err = io.ReadFull(r.reader, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
