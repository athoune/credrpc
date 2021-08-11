package main

import (
	"encoding/gob"
	"errors"
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
	s := server.NewServer(func(i *gob.Decoder, u *syscall.Ucred) (interface{}, error) {
		if u.Uid == 0 {
			return nil, errors.New("root is not allowed.")
		}
		var data string
		err := i.Decode(&data)
		if err != nil {
			return nil, err
		}
		fmt.Println("msg", data)
		fmt.Println("user", u)
		if data == "pam" {
			return nil, errors.New("I don't like Pam")
		}

		return fmt.Sprintf("Hello %s", data), nil
	})

	err := s.ListenAndServe(listen)

	if err != nil {
		log.Fatal(err)
	}
}
