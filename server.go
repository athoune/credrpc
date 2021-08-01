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
	l, err := net.Listen("unix", listen)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		c, ok := fd.(*net.UnixConn)
		if !ok {
			log.Fatal("Not a UNIX socket")
		}
		f, err := c.File()
		if err != nil {
			log.Fatal("socket file error:", err)
		}
		//s, _ := f.Stat()

		err = syscall.SetsockoptInt(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_PASSCRED, 1)
		if err != nil {
			panic(err)
		}

		fmt.Println(f.Name())

		go echoServer(fd)
	}
}
