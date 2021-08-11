package client

import (
	"encoding/gob"
	"net"
	"os"
	"syscall"
)

type Client struct {
	conn *net.UnixConn
	enc  *gob.Encoder
	dec  *gob.Decoder
	oob  []byte
}

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

func (c *Client) Do(input interface{}, output interface{}) error {
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
