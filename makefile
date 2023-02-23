.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY: fmt

staticcheck: fmt
	staticcheck ./...
.PHONY: staticcheck

vet: staticcheck
	go vet ./...
.PHONY: vet

build: vet
	mkdir -p bin
	go build -o ./bin ./...
.PHONY: build.DEFAULT_GOAL := build
