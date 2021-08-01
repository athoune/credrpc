package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		println("Client got:", string(buf[0:n]))
	}
}

func main() {
	server := os.Getenv("SERVER")
	if server == "" {
		server = "/tmp/echo.sock"
	}
	c, err := net.Dial("unix", server)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	cc, ok := c.(*net.UnixConn)
	if !ok {
		panic("not a unix socket")
	}
	/*
		f, err := cc.File()
		if err != nil {
			panic(err)
		}

			err = syscall.SetsockoptInt(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_PASSCRED, 1)
			if err != nil {
				panic(err)
			}
	*/

	var ucred syscall.Ucred
	ucred.Pid = int32(os.Getpid())
	ucred.Uid = uint32(os.Getuid())
	ucred.Gid = uint32(os.Getgid())

	oob := syscall.UnixCredentials(&ucred)

	go reader(c)
	for {
		n, oobn, err := cc.WriteMsgUnix([]byte("hi"), oob, nil)
		fmt.Println(n, oobn)
		if err != nil {
			log.Fatal("write error:", err)
			break
		}
		time.Sleep(1e9)
	}
}
