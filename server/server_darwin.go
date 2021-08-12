package server

import (
	"net"
	"syscall"
)

// FIXME

func CredLen() int {
	return 4
}

func PrepareSocket(c *net.UnixConn) error {
	return nil
}

func SocketControlMessage2Cred(scm []syscall.SocketControlMessage) (*Cred, error) {
	return nil, nil
}
