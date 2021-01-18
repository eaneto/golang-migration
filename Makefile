build: clean
	mkdir -p bin
	go build  -o ./bin/ -v ./cmd/grotto

clean:
	rm -rf bin

test: build
	go test -v ./...

codecov: test
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic
