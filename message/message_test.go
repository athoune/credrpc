package message

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	c := &bytes.Buffer{}
	w := NewWriter(c)
	n, err := w.Write([]byte("Beuha"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	r := NewReader(c)
	resp, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, []byte("Beuha"), resp, string(resp))
}
