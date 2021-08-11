package client

import (
	"encoding/gob"
	"errors"
	"net"
	"os"
	"syscall"
)

// Client talks to the server, with UNIX credential and gob encoding
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
func (c *Client) Call(input interface{}, output interface{}) error {
	cc, err := net.Dial("unix", c.path)
	if err != nil {
		return err
	}
	defer cc.Close()
	conn := cc.(*net.UnixConn)
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	oob := syscall.UnixCredentials(&syscall.Ucred{
		Pid: int32(os.Getpid()),
		Uid: uint32(os.Getuid()),
		Gid: uint32(os.Getgid()),
	})
	_, _, err = conn.WriteMsgUnix(nil, oob, nil)
	if err != nil {
		return err
	}
	err = enc.Encode(input)
	if err != nil {
		return err
	}
	var errRpc string
	err = dec.Decode(&errRpc)
	if err != nil {
		return err
	}
	if len(errRpc) != 0 {
		return errors.New(errRpc)
	}
	err = dec.Decode(output)
	if err != nil {
		return err
	}
	return nil
}
