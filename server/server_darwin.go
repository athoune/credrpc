package server

import (
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func CredLen() int {
	return 4
}

func PrepareSocket(l *net.UnixListener) error {
	f, err := l.File()
	if err != nil {
		return err
	}
	_, err = unix.GetsockoptXucred(int(f.Fd()), syscall.SOL_SOCKET, unix.LOCAL_PEERCRED)
	return err
}

func SocketControlMessage2Cred(scm []syscall.SocketControlMessage) (*Cred, error) {
	return nil, nil
}

func ActivationListener() (net.Listener, error) {
	return nil, nil
}
