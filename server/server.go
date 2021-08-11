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
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		f, err := conn.(*net.UnixConn).File()
		if err != nil {
			return err
		}
		err = syscall.SetsockoptInt(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_PASSCRED, 1)
		if err != nil {
			return err
		}
		enc := gob.NewEncoder(conn)
		dec := gob.NewDecoder(conn)
		go func() {
			for {
				oob2 := make([]byte, len(syscall.UnixCredentials(&syscall.Ucred{})))
				_, _, flags, _, err := conn.(*net.UnixConn).ReadMsgUnix(nil, oob2)
				if err != nil {
					log.Print(err)
					conn.Close()
				}
				if flags != 0 {
					log.Fatal("Strange flags", flags)
				}
				scm, err := syscall.ParseSocketControlMessage(oob2)
				if err != nil {
					log.Print(err)
					conn.Close()
				}
				newUcred, err := syscall.ParseUnixCredentials(&scm[0])
				if err != nil {
					log.Print(err)
					conn.Close()
				}
				err = s.handler(dec, enc, newUcred)
				if err != nil {
					log.Print(err)
					conn.Close()
				}
			}
		}()
	}
	return nil
}
