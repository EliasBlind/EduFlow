# Project Variables
MODULE = github.com/EliasBlind/EduFlow
PROTO_SRC = protos/journal/v1
GEN_OUT = pkg/protos/gen
REPORT_PATH = docs/reports
COVER_OUT = $(REPORT_PATH)/coverage.out

.PHONY: all gen clean rebuild test cover cover-html build help install-deps

all: gen build

## install-deps: Install protoc dependencies for Go
install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

## gen: Generate Go code from journal.proto
gen:
	@mkdir -p $(GEN_OUT)
	protoc --proto_path=protos \
		--go_out=$(GEN_OUT) --go_opt=module=$(MODULE)/pkg/protos/gen \
		--go-grpc_out=$(GEN_OUT) --go-grpc_opt=module=$(MODULE)/pkg/protos/gen \
		$(PROTO_SRC)/*.proto


## build: Compile the server binary
build:
	go build -o bin/server ./cmd/server/main.go

## clean: Remove generated files and binaries
clean:
	rm -rf $(GEN_OUT)/*
	rm -rf bin/
	rm -rf $(REPORT_PATH)

rebuild: clean
	go mod tidy
	$(MAKE) gen
	$(MAKE) build

## test: Run all tests in the project
test:
	go test -v ./...

## cover: Run tests and show coverage percentage in terminal
cover:
	go test -coverprofile=$(COVER_OUT) ./...
	go tool cover -func=$(COVER_OUT)

## cover-html: Generate HTML coverage report and open it in browser
cover-std-html:
	@mkdir -p $(REPORT_PATH)
	go test -coverprofile=$(COVER_OUT) ./...
	go tool cover -html=$(COVER_OUT) -o $(REPORT_PATH)/coverage.html
	xdg-open $(REPORT_PATH)/coverage.html

## cover-nice: Generate a beautiful coverage report
cover-html:
	@mkdir -p $(REPORT_PATH)
	gocov test ./pkg/utils/convert/... | gocov-html > $(REPORT_PATH)/pretty_coverage.html
	xdg-open $(REPORT_PATH)/pretty_coverage.html

## help: Show available commands
help:
	# TODO: Написать help
