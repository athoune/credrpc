package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/factorysh/chownme/client"
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

	cli := client.New(c.(*net.UnixConn))

	for {
		var data string
		err = cli.Call("World", &data)
		if err != nil {
			log.Fatal("write error:", err)
		}
		println("Client got:", data)
		time.Sleep(1e9)
	}
}
