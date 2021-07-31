package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
)

func echoServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

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
	l, err := net.Listen("unix", listen)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	for {
		fd, err := l.Accept()
		c, _ := fd.(*net.UnixConn)
		f, _ := c.File()
		//s, _ := f.Stat()
		err = syscall.SetsockoptInt(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_PASSCRED, 1)
		fmt.Println(f.Name())
		if err != nil {
			log.Fatal("accept error:", err)
		}

		go echoServer(fd)
	}
}
