package main

import (
	"log"
	"os"
	"time"

	"github.com/athoune/credrpc/client"
)

func main() {
	server := os.Getenv("SERVER")
	if server == "" {
		server = "/tmp/echo.sock"
	}

	cli := client.New(server)

	for _, name := range []string{"Pim", "Pam", "Poum"} {
		data, err := cli.Call([]byte(name))
		if err != nil {
			log.Print("Call error:", err)
		}
		println("Client got:", string(data))
		time.Sleep(1e9)
	}
}
