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
			defer conn.Close()
			oob2 := make([]byte, len(syscall.UnixCredentials(&syscall.Ucred{})))
			n, _, flags, _, err := conn.(*net.UnixConn).ReadMsgUnix(nil, oob2)
			if err != nil {
				if n == 0 { // conn seems to be closed
					log.Print("Closed UNIX socket : ", f.Name())
				} else {
					log.Print("Can't read header : ", err)
				}
				return
			}
			if flags != 0 {
				log.Fatal("Strange flags ", flags)
			}
			scm, err := syscall.ParseSocketControlMessage(oob2)
			if err != nil {
				log.Print("Can't parse socket control message : ", err)
				return
			}
			newUcred, err := syscall.ParseUnixCredentials(&scm[0])
			if err != nil {
				log.Print("Can't parse UNIX credential : ", err)
				return
			}
			err = s.handler(dec, enc, newUcred)
			if err != nil {
				log.Print("Error Handler : ", err)
				enc.Encode(err.Error())
				return
			} else {
				err = enc.Encode("")
				if err != nil {
					log.Print("Error while returnging response : ", err)
				}
			}
		}()
	}
	return nil
}
