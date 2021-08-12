package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/athoune/credrpc/server"
)

func main() {
	listener, err := server.ActivationListener()
	if err != nil {
		log.Fatal(err)
	}
	if listener == nil {
		listen := os.Getenv("LISTEN")
		if listen == "" {
			listen = "/tmp/echo.sock"
		}
		listener, err = net.Listen("unix", listen)
		if err != nil {
			log.Fatal(err)
		}
	}
	s := server.NewServer(func(i []byte, u *server.Cred) ([]byte, error) {
		if u.Uid == 0 {
			return nil, errors.New("root is not allowed")
		}
		data := string(i)
		fmt.Println("msg", data)
		fmt.Println("user", u)
		if data == "Pam" {
			return nil, errors.New("I don't like Pam")
		}

		return []byte(fmt.Sprintf("Hello %s", data)), nil
	})

	err = s.Serve(listener)

	if err != nil {
		log.Fatal(err)
	}
}
