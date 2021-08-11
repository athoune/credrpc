package main

import (
	"encoding/gob"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

func main() {
	server := os.Getenv("SERVER")
	if server == "" {
		server = "/tmp/echo.sock"
	}
	c, err := net.Dial("unix", server)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	enc := gob.NewEncoder(c)
	dec := gob.NewDecoder(c)

	var ucred syscall.Ucred
	ucred.Pid = int32(os.Getpid())
	ucred.Uid = uint32(os.Getuid())
	ucred.Gid = uint32(os.Getgid())

	oob := syscall.UnixCredentials(&ucred)

	for {
		_, _, err := c.(*net.UnixConn).WriteMsgUnix(nil, oob, nil)
		if err != nil {
			log.Fatal("write error:", err)
		}
		err = enc.Encode("World")
		if err != nil {
			log.Fatal("write error:", err)
		}
		var data string
		err = dec.Decode(&data)
		if err != nil {
			log.Fatal(err)
		}
		println("Client got:", data)
		time.Sleep(1e9)
	}
}
