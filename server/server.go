package server

import (
	"log"
	"net"
	"syscall"
	"time"

	"github.com/factorysh/chownme/protocol"
)

type Cred struct {
	Pid int32
	Uid uint32
	Gid uint32
}

type Handler func(i []byte, c *Cred) ([]byte, error)

type Server struct {
	handler Handler
}

func NewServer(handler Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) ListenAndServe(path string) error {
	listener, err := net.Listen("unix", path)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		conn.SetDeadline(time.Now().Add(time.Second))

		err = PrepareSocket(conn.(*net.UnixConn))
		if err != nil {
			return err
		}

		go func(conn net.Conn) {
			defer conn.Close()
			oob2 := make([]byte, CredLen())
			buff := make([]byte, 2*1024) // 2k should be enough
			n, _, flags, _, err := conn.(*net.UnixConn).ReadMsgUnix(buff, oob2)
			if err != nil {
				if n == 0 { // conn seems to be closed
					log.Print("Closed UNIX socket : ", conn)
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
			newUcred, err := SocketControlMessage2Cred(scm)
			if err != nil {
				log.Print("Can't parse UNIX credential : ", err)
				return
			}
			input, err := protocol.Read(buff[:n], conn)
			if err != nil {
				log.Print("Can't read input : ", err)
				return
			}

			resp, err := s.handler(input, newUcred)
			if err != nil {
				log.Print("Error Handler : ", err)
				err = protocol.Write(conn, []byte(err.Error()))
				if err != nil {
					log.Print("Can't write error : ", err)
					return
				}
				// don't bother to send nil response, connection will be closed
			} else {
				err = protocol.Write(conn, []byte{})
				if err != nil {
					log.Print("Error while returnging empty error : ", err)
				}
				err = protocol.Write(conn, resp)
				if err != nil {
					log.Print("Error while returnging response : ", err)
				}
			}
		}(conn)
	}
	return nil
}
