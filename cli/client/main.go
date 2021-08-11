package main

import (
	"log"
	"os"
	"time"

	"github.com/factorysh/chownme/client"
)

func main() {
	server := os.Getenv("SERVER")
	if server == "" {
		server = "/tmp/echo.sock"
	}

	cli := client.New(server)
	var err error

	for _, name := range []string{"pim", "pam", "poum"} {
		var data string
		err = cli.Call(name, &data)
		if err != nil {
			log.Fatal("write error:", err)
		}
		println("Client got:", data)
		time.Sleep(1e9)
	}
}
