# Project Variables
MODULE = github.com/EliasBlind/EduFlow
PROTO_SRC = protos/journal/v1
GEN_OUT = pkg/protos/gen

.PHONY: all gen clean rebuild test build help install-deps

all: gen build

## install-deps: Install protoc dependencies for Go
install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

## gen: Generate Go code from journal.proto
gen:
	@mkdir -p $(GEN_OUT)
	protoc --proto_path=protos \
		--go_out=$(GEN_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)/*.proto

## build: Compile the server binary
build:
	go build -o bin/server ./cmd/server/main.go

## clean: Remove generated files and binaries
clean:
	rm -rf $(GEN_OUT)/*
	rm -rf bin/

rebuild: clean
	go mod tidy
	$(MAKE) gen
	$(MAKE) build

## help: Show available commands
help:
	# TODO: Написать help
