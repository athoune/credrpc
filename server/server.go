package server

import (
	"io"
	"log"
	"net"
	"syscall"
)

type Handler func(i []byte, o io.Writer, u *syscall.Ucred) error

type Server struct {
	listener net.Listener
	handler  Handler
}

func NewServer(handler Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) ListenAndServe(path string) error {
	var err error
	s.listener, err = net.Listen("unix", path)
	if err != nil {
		return err
	}
	for {
		fd, err := s.listener.Accept()
		if err != nil {
			return err
		}
		f, err := fd.(*net.UnixConn).File()
		if err != nil {
			return err
		}
		err = syscall.SetsockoptInt(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_PASSCRED, 1)
		if err != nil {
			return err
		}
		go func() {
			for {
				buf := make([]byte, 512)
				oob2 := make([]byte, 10*24)
				nr, oobn2, flags, _, err := fd.(*net.UnixConn).ReadMsgUnix(buf, oob2)
				if err != nil {
					log.Fatal(err)
				}
				if flags != 0 {
					log.Fatal("Strange flags", flags)
				}
				oob2 = oob2[:oobn2]
				scm, err := syscall.ParseSocketControlMessage(oob2)
				if err != nil {
					log.Fatal(err)
				}
				newUcred, err := syscall.ParseUnixCredentials(&scm[0])
				if err != nil {
					log.Fatal(err)
				}
				data := buf[0:nr]
				err = s.handler(data, fd, newUcred)
				if err != nil {
					log.Fatal(err)
				}
			}
		}()
	}
	return nil
}
