.DEFAULT_GOAL := build

ifeq ($(OS),Windows_NT)
define MD
	if not exist $(1) mkdir $(1)
endef
else
define MD
	mkdir -p $(1)
endef
endif

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
	$(call MD, bin)
	go build -o ./bin ./...
.PHONY: build.DEFAULT_GOAL := build

generate:
	go get github.com/99designs/gqlgen@v0.17.31 && go generate ./...