package client

import (
	"encoding/gob"
	"net"
	"os"
	"syscall"
)

// Client talks to the server, with UNIX credential and gob encoding
type Client struct {
	conn *net.UnixConn
	enc  *gob.Encoder
	dec  *gob.Decoder
	oob  []byte
}

// New Client, don't forget to close the connection with a defer.
func New(conn *net.UnixConn) *Client {
	return &Client{
		conn: conn,
		enc:  gob.NewEncoder(conn),
		dec:  gob.NewDecoder(conn),
		oob: syscall.UnixCredentials(&syscall.Ucred{
			Pid: int32(os.Getpid()),
			Uid: uint32(os.Getuid()),
			Gid: uint32(os.Getgid()),
		}),
	}
}

// Call the server with an input and an output pointer for the answer.
func (c *Client) Call(input interface{}, output interface{}) error {
	_, _, err := c.conn.WriteMsgUnix(nil, c.oob, nil)
	if err != nil {
		c.conn.Close()
		return err
	}
	err = c.enc.Encode(input)
	if err != nil {
		c.conn.Close()
		return err
	}
	err = c.dec.Decode(output)
	if err != nil {
		c.conn.Close()
		return err
	}
	return nil
}
