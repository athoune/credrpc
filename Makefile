.PHONY: client server darwin

build: client server

server: bin
	GOOS=linux go build -o bin/chownmed ./cli/server/main.go

client: bin
	GOOS=linux go build -o bin/chownme ./cli/client/main.go

bin:
	mkdir -p bin

darwin/bin:
	mkdir -p darwin/bin

darwin: darwin/bin
	go build -o darwin/bin/chownmed ./cli/server/main.go
	go build -o darwin/bin/chownme ./cli/client/main.go

run/chownme:
	mkdir -p run/chownme

up: run/chownme
	rm -rf run/chownme/sock
	docker-compose up --exit-code-from client

test:
	go test -cover \
		github.com/athoune/credrpc/protocol
