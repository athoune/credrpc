package server

import (
	"fmt"
	"log"
	"net"
	"syscall"
	"time"

	"github.com/athoune/credrpc/protocol"
)

type Cred struct {
	Pid int32
	Uid uint32
	Gid uint32
}

type Handler func(i []byte, c *Cred) ([]byte, error)
type Logger func(error)

type Server struct {
	handler Handler
	logger  Logger
}

func NewServer(handler Handler) *Server {
	return &Server{
		handler: handler,
		logger: func(e error) {
			log.Print(e)
		},
	}
}

func (s *Server) Serve(listener net.Listener) error {
	defer listener.Close()
	u := listener.(*net.UnixListener)
	err := PrepareSocket(u)
	if err != nil {
		return err
	}
	listener.(*net.UnixListener).SetUnlinkOnClose(true)
	f, err := u.File()
	if err != nil {
		return err
	}
	err = f.Chmod(0770)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func(conn net.Conn) {
			defer conn.Close()
			err := conn.SetDeadline(time.Now().Add(3 * time.Minute))
			if err != nil {
				s.logger(err)
				return
			}
			err = s.handle(conn)
			if err != nil {
				s.logger(err)
			}
		}(conn)
	}
	return nil
}

func (s *Server) handle(conn net.Conn) error {
	defer conn.Close()
	oob2 := make([]byte, CredLen())
	buff := make([]byte, 2*1024) // 2k should be enough
	n, _, flags, _, err := conn.(*net.UnixConn).ReadMsgUnix(buff, oob2)
	if err != nil {
		if n == 0 { // conn seems to be closed
			return fmt.Errorf("closed UNIX socket : %v", conn)
		} else {
			return fmt.Errorf("can't read header : %v", err)
		}
	}
	if flags != 0 {
		return fmt.Errorf("strange flags %v", flags)
	}
	scm, err := syscall.ParseSocketControlMessage(oob2)
	if err != nil {
		return fmt.Errorf("can't parse socket control message : %v", err)
	}
	newUcred, err := SocketControlMessage2Cred(scm)
	if err != nil {
		return fmt.Errorf("can't parse UNIX credential : %v", err)
	}
	input, err := protocol.Read(buff[:n], conn)
	if err != nil {
		return fmt.Errorf("can't read input : %v", err)
	}

	resp, err := s.handler(input, newUcred)
	if err != nil {
		s.logger(err)
		err = protocol.Write(conn, []byte(err.Error()))
		if err != nil {
			return fmt.Errorf("can't write error : %v", err)
		}
		// don't bother to send nil response, connection will be closed
	} else {
		err = protocol.Write(conn, []byte{}, resp)
		if err != nil {
			return fmt.Errorf("error while returnging response : %v", err)
		}
	}
	return nil
}
