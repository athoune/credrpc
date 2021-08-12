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

	for _, name := range []string{"pim", "pam", "poum"} {
		data, err := cli.Call([]byte(name))
		if err != nil {
			log.Print("Call error:", err)
		}
		println("Client got:", string(data))
		time.Sleep(1e9)
	}
}
