# Project Settings
MODULE              = github.com/EliasBlind/EduFlow

# Proto Sources
JOURNAL_PROTO_SRC   = protos/journal/v1
SSO_PROTO_SRC       = protos/sso/v1

# Generation and Reports
COVER_OUT           = $(REPORT_PATH)/coverage.out
GEN_OUT             = pkg/protos/gen
REPORT_PATH         = docs/reports

# Binaries
ENVGEN              = .bin/envgen
JOURNAL_SERVICE     = .bin/journal_service
SSO_SERVICE         = .bin/sso_service

# Services
SSO_DIR				= internal/sso_service
JOURNAL_DIR			= internal/journal_service

DB_URL="postgres://Elias:2795749b040201cde46538c3b74b4d97@127.0.0.1:5432/sso_db?sslmode=disable"

.PHONY: all gen clean rebuild test cover cover-html build help install-deps

all: gen build

## install-deps: Install protoc dependencies for Go
install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

## build: Compile the server binary
journal-build:
	go build -o $(JOURNAL_SERVICE) ./cmd/journal_service/server/main.go

journal-run:
	go run ./cmd/journal_service/server/main.go

sso-build:
	go build -o $(JOURNAL_SERVICE) ./cmd/sso_service/server/main.go

sso-run:
	-go run ./cmd/sso_service/server/main.go

## clean: Remove generated files and binaries
clean:
	rm -rf bin/
	rm -rf $(GEN_OUT)/*
	rm -rf $(REPORT_PATH)
	rm -rf $(SSO_DIR)/storage/sqlgen
	rm -rf $(JOURNAL_DIR)/storage/db


rebuild: clean
	go mod tidy
	$(MAKE) gen
	$(MAKE) build

## gen: Generate Go code

journal-gen-proto:
	mkdir -p $(GEN_OUT)
	protoc --proto_path=protos \
		--go_out=$(GEN_OUT) --go_opt=module=$(MODULE)/pkg/protos/gen \
		--go-grpc_out=$(GEN_OUT) --go-grpc_opt=module=$(MODULE)/pkg/protos/gen \
		$(JOURNAL_PROTO_SRC)/*.proto

journal-gen-sqlc:
	@echo "Generating SQLC for Journal Service..."
	cd $(JOURNAL_DIR) && sqlc generate

sso-gen-proto:
	mkdir -p $(GEN_OUT)
	protoc --proto_path=protos \
		--go_out=$(GEN_OUT) --go_opt=module=$(MODULE)/pkg/protos/gen \
		--go-grpc_out=$(GEN_OUT) --go-grpc_opt=module=$(MODULE)/pkg/protos/gen \
		$(SSO_PROTO_SRC)/*.proto

sso-gen-sqlc:
	@echo "Generating SQLC for SSO Service..."
	cd $(SSO_DIR) && sqlc generate

# Env data generate

key-gen: build-envgen
	$(eval NEW_VAL := $(shell if [ -z "$(Key)" ]; then openssl rand -hex 32; else echo "$(Key)"; fi))
	@./.bin/envgen -env="configs/journal_service/journal.env" -key="SECRET_KEY" -sed="$(NEW_VAL)"
	@./.bin/envgen -env="configs/sso_service/sso.env" -key="SECRET_KEY" -sed="$(NEW_VAL)"
	@echo "The SAME key has been updated in both env files"


storage-passwd-gen: build-envgen
	@$(if $(Env),,$(eval Env=.env))
	@$(ENVGEN) -env="$(Env)" -key="STORAGE_PASSWORD" -sed="$(Passwd)" -gen="db"
	@echo "Password updated in $(Env)"

sso-gen: update_proto build-envgen
	@$(ENVGEN) -env="configs/sso_service/sso.env" -key="$(KEY)" -sed="$(SED)" -gen="$(GEN)"
	@echo "SSO variable $(KEY) updated"

journal-gen: update_proto build-envgen
	@$(ENVGEN) -env="configs/journal_service/journal.env" -key="$(KEY)" -sed="$(SED)" -gen="$(GEN)"
	@echo "Journal variable $(KEY) updated"

update_proto:
	git submodule update --remote --merge protos

## test: Run all tests in the project
test:
	go test -v ./...

## cover: Run tests and show coverage percentage in terminal
cover: gen
	go test -coverprofile=$(COVER_OUT) ./...
	go tool cover -func=$(COVER_OUT)

## cover-html: Generate HTML coverage report and open it in browser
cover-std-html: gen
	mkdir -p $(REPORT_PATH)
	go test -coverprofile=$(COVER_OUT) ./...
	go tool cover -html=$(COVER_OUT) -o $(REPORT_PATH)/coverage.html
	xdg-open $(REPORT_PATH)/coverage.html

## cover-nice: Generate a beautiful coverage report
cover-html: gen
	mkdir -p $(REPORT_PATH)
	gocov test ./pkg/utils/convert/... | gocov-html > $(REPORT_PATH)/pretty_coverage.html
	xdg-open $(REPORT_PATH)/pretty_coverage.html

build-envgen:
	mkdir -p .bin
	go build -o $(ENVGEN) cmd/envgen/main.go

migrate-up:
	goose -dir internal/sso_service/sql/schema postgres $(DB_URL) up

migrate-down:
	goose -dir internal/sso_service/sql/schema postgres $(DB_URL) down

## help: Show available commands
help:
	# TODO: Написать help
