package message

import (
	"encoding/binary"
	"io"
)

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
