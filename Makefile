all:
	rm -rf bin
	mkdir -p bin
	go build -o bin/

test:
	go test -v ./...
