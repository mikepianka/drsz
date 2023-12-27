.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY: fmt

vet:
	go vet ./...
.PHONY: vet

lint: fmt vet
	staticcheck ./...
.PHONY: lint

test: lint
	go test ./...
.PHONY: test

clean:
	rm -rf bin/
.PHONY: clean

build-linux: test clean
	GOOS=linux GOARCH=amd64 go build -o bin/linux/ ./...
.PHONY: build-linux

build-mac: test clean
	GOOS=darwin GOARCH=arm64 go build -o bin/mac/ ./...
.PHONY: build-mac

build-windows: test clean
	GOOS=windows GOARCH=amd64 go build -o bin/windows/ ./...
.PHONY: build-windows

build: build-windows build-mac build-linux
.PHONY: build