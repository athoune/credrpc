package client

import (
	"encoding/binary"
	"errors"
	"net"
	"os"
	"syscall"

	"github.com/factorysh/chownme/protocol"
)

// Client talks to the server, with UNIX credential
type Client struct {
	path string
}

// New Client, don't forget to close the connection with a defer.
func New(path string) *Client {
	return &Client{
		path: path,
	}
}

// Call the server with an input and an output pointer for the answer.
func (c *Client) Call(input []byte) ([]byte, error) {
	cc, err := net.Dial("unix", c.path)
	if err != nil {
		return nil, err
	}
	defer cc.Close()
	conn := cc.(*net.UnixConn)
	oob := syscall.UnixCredentials(&syscall.Ucred{
		Pid: int32(os.Getpid()),
		Uid: uint32(os.Getuid()),
		Gid: uint32(os.Getgid()),
	})
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(input)))
	_, _, err = conn.WriteMsgUnix(size, oob, nil)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(input)
	if err != nil {
		return nil, err
	}
	errRpc, err := protocol.Read([]byte{}, conn)
	if err != nil {
		return nil, err
	}
	if len(errRpc) != 0 {
		return nil, errors.New(string(errRpc))
	}
	output, err := protocol.Read([]byte{}, conn)
	if err != nil {
		return nil, err
	}
	return output, nil
}
