package protocol

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocol(t *testing.T) {
	b := &bytes.Buffer{}
	err := Write(b, []byte("beuha"))
	assert.NoError(t, err)
	r, err := Read([]byte{}, b)
	assert.NoError(t, err)
	assert.Equal(t, []byte("beuha"), r)
}

func TestProtocolWithBuff(t *testing.T) {
	b := bytes.NewBuffer([]byte("beuha"))
	s := make([]byte, 4)
	binary.BigEndian.PutUint32(s, 5)
	r, err := Read(s, b)
	assert.NoError(t, err)
	assert.Equal(t, []byte("beuha"), r)
}

func TestProtocolWithFatBuffer(t *testing.T) {
	b := make([]byte, 2048)
	binary.BigEndian.PutUint32(b[:4], 5)
	copy(b[4:], []byte("Beuha"))
	r, err := Read(b, &bytes.Buffer{})
	assert.NoError(t, err)
	assert.Equal(t, []byte("Beuha"), r)
}
