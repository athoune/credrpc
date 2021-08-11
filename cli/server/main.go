package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/factorysh/chownme/server"
)

func main() {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "/tmp/echo.sock"
	}
	s := server.NewServer(func(i *gob.Decoder, o *gob.Encoder, u *syscall.Ucred) error {
		var data string
		err := i.Decode(&data)
		if err != nil {
			return err
		}
		fmt.Println("msg", data)
		fmt.Println("user", u)
		return o.Encode(fmt.Sprintf("Hello %s", data))
	})

	err := s.ListenAndServe(listen)

	if err != nil {
		log.Fatal(err)
	}
}
