package server

// The patch is merged https://go-review.googlesource.com/c/sys/+/292330/
// Not yet available in go 1.16

import (
	"net"
	"syscall"
)

func CredLen() int {
	return 4
}

func PrepareSocket(c *net.UnixListener) error {
	return nil
}

func SocketControlMessage2Cred(scm []syscall.SocketControlMessage) (*Cred, error) {
	return nil, nil
}

func ActivationListener() (net.Listener, error) {
	return nil, nil
}
