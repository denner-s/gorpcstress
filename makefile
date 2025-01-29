BINARY=gorpcstress
VERSION=1.0.0

build:
	go build -o bin/$(BINARY) ./cmd/gorpcstress

test:
	go test -v ./...

lint:
	golangci-lint run

bench:
	go test -bench=. -benchmem ./internal/...

clean:
	rm -rf bin/

release:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY)-linux-amd64 ./cmd/gorpcstress
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY)-darwin-amd64 ./cmd/gorpcstress

.PHONY: build test lint bench clean release