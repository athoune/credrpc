package server

import (
	"encoding/gob"
	"log"
	"net"
	"syscall"
)

type Handler func(i *gob.Decoder, o *gob.Encoder, u *syscall.Ucred) error

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
		enc := gob.NewEncoder(fd)
		dec := gob.NewDecoder(fd)
		go func() {
			for {
				oob2 := make([]byte, len(syscall.UnixCredentials(&syscall.Ucred{})))
				_, _, flags, _, err := fd.(*net.UnixConn).ReadMsgUnix(nil, oob2)
				if err != nil {
					log.Print(err)
					fd.Close()
				}
				if flags != 0 {
					log.Fatal("Strange flags", flags)
				}
				scm, err := syscall.ParseSocketControlMessage(oob2)
				if err != nil {
					log.Print(err)
					fd.Close()
				}
				newUcred, err := syscall.ParseUnixCredentials(&scm[0])
				if err != nil {
					log.Print(err)
					fd.Close()
				}
				err = s.handler(dec, enc, newUcred)
				if err != nil {
					log.Print(err)
					fd.Close()
				}
			}
		}()
	}
	return nil
}
