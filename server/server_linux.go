package server

import (
	"net"
	"syscall"
)

func CredLen() int {
	return len(syscall.UnixCredentials(&syscall.Ucred{}))
}

func PrepareSocket(c *net.UnixConn) error {
	f, err := c.File()
	if err != nil {
		return err
	}

	// Please, pass credential on the socket
	return syscall.SetsockoptInt(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_PASSCRED, 1)
}

func SocketControlMessage2Cred(scm []syscall.SocketControlMessage) (*Cred, error) {
	newUcred, err := syscall.ParseUnixCredentials(&scm[0])
	if err != nil {
		return nil, err
	}
	return &Cred{
		Pid: newUcred.Pid,
		Uid: newUcred.Uid,
		Gid: newUcred.Gid,
	}, nil
}
