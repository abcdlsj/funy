build:
	rm -rf bin
	
	go build -o bin/funy cmd/funy/cli.go
	go build -o bin/funyd cmd/funyd/server.go