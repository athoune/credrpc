package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"

	"github.com/factorysh/chownme/server"
)

func echoServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		oob2 := make([]byte, 10*24)
		nr, oobn2, flags, _, err := c.(*net.UnixConn).ReadMsgUnix(buf, oob2)
		if err != nil {
			return
		}

		fmt.Println("flags", flags)

		oob2 = oob2[:oobn2]
		scm, err := syscall.ParseSocketControlMessage(oob2)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("scm", scm)
		newUcred, err := syscall.ParseUnixCredentials(&scm[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("ucred", newUcred)

		data := buf[0:nr]
		println("Server got:", string(data))
		_, err = c.Write(data)
		if err != nil {
			log.Fatal("Write: ", err)
		}
	}
}

func main() {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "/tmp/echo.sock"
	}
	s := server.NewServer(func(i []byte, o io.Writer, u *syscall.Ucred) error {
		fmt.Println("msg", i)
		o.Write(i)
		fmt.Println("user", u)
		return nil
	})

	err := s.ListenAndServe(listen)

	if err != nil {
		log.Fatal(err)
	}
}
