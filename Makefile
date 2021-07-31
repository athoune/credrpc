.PHONY: client server

build: client server

server: bin
	GOOS=linux go build -o bin/chownmed ./server.go

client: bin
	GOOS=linux go build -o bin/chownme ./client.go

bin:
	mkdir -p bin

run/chownme:
	mkdir -p run/chownme

up: run/chownme
	rm -rf run/chownme/sock
	docker-compose up
