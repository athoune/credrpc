package main

import (
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
	s := server.NewServer(func(i []byte, u *syscall.Ucred) ([]byte, error) {
		if u.Uid == 0 {
			return nil, errors.New("root is not allowed.")
		}
		data := string(i)
		fmt.Println("msg", data)
		fmt.Println("user", u)
		if data == "pam" {
			return nil, errors.New("I don't like Pam")
		}

		return []byte(fmt.Sprintf("Hello %s", data)), nil
	})

	err := s.ListenAndServe(listen)

	if err != nil {
		log.Fatal(err)
	}
}
