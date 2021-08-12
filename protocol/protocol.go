package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
)

func Write(w io.Writer, data []byte) error {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, uint32(len(data)))
	_, err := w.Write(buff)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func Read(stack []byte, r io.Reader) ([]byte, error) {
	buffer := bytes.NewBuffer(stack)
	var err error
	if buffer.Len() < 4 {
		_, err = io.CopyN(buffer, r, int64(4-buffer.Len()))
		if err != nil {
			return nil, err
		}
	}
	size := binary.BigEndian.Uint32(buffer.Bytes()[0:4])
	if size == 0 {
		return []byte{}, nil
	}
	if size > uint32(buffer.Len()-4) { // the buffer is not complete, lets read
		_, err = io.CopyN(buffer, r, int64(size-uint32(buffer.Len()-4)))
		if err != nil {
			return nil, err
		}
	}
	return buffer.Bytes()[4 : size+4], nil
}
