all:
	rm -rf bin
	mkdir -p bin
	go build -o bin/

test:
	go test -v ./...

codecov:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic
