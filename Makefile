# Project Settings
MODULE              = github.com/EliasBlind/EduFlow

# Proto Sources
PROTO_ROOT          = protos
GOOGLEAPIS_PATH     = third_party/googleapis

# Generation and Reports
COVER_OUT           = $(REPORT_PATH)/coverage.out
GEN_OUT             = pkg/protos/gen
REPORT_PATH         = docs/reports

# Binaries
ENVGEN              = .bin/envgen
JOURNAL_SERVICE_BIN = .bin/journal_service
SSO_SERVICE_BIN     = .bin/sso_service

# Services
SSO_DIR             = internal/sso_service
JOURNAL_DIR         = internal/journal_service

DB_URL="postgres://Elias:2795749b040201cde46538c3b74b4d97@127.0.0.1:5432/edu_db?sslmode=disable"


.PHONY: all gen clean rebuild test cover cover-html run build build-all journal-build sso-build \
        journal-run sso-run help migrate-up migrate-down dependencies_install \
        proto-deps proto-gen proto-descriptor journal-gen-sqlc sso-gen-sqlc \
        init-submodules update_proto

all: gen build

## build: Compile both server binaries
build: dependencies_install gen journal-build sso-build

run: build
	docker compose up -d
	$(MAKE) migrate-up

## journal-build: Compile the journal server binary
journal-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(JOURNAL_SERVICE_BIN) ./cmd/journal_service/server/main.go 

## sso-build: Compile the SSO server binary
sso-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(SSO_SERVICE_BIN) ./cmd/sso_service/server/main.go 

dependencies_install:
	cd ./cmd/dependencies/ && bash install.sh

## init-submodules: Initialize and update git submodules (if any)
init-submodules:
	@echo "Initializing submodules..."
	git submodule update --init --recursive

## update_proto: Update protos submodule to the latest remote commit (develop branch)
update_proto: init-submodules
	@echo "Updating protos submodule..."
	git submodule update --remote --merge protos

## clean: Remove generated files and binaries
clean:
	docker compose down
	rm -rf .bin/ bin/
	rm -rf $(GEN_OUT)/*
	rm -rf $(REPORT_PATH)
	rm -rf $(SSO_DIR)/storage/sqlgen
	rm -rf $(JOURNAL_DIR)/storage/sqlgen

rebuild: clean
	go mod tidy
	$(MAKE) gen
	$(MAKE) build

## gen: Generate Go code (protobuf + SQLC)
gen: proto-deps proto-gen proto-descriptor journal-gen-sqlc sso-gen-sqlc

## proto-deps: Ensure googleapis are available (clone if missing, skip update if present)
proto-deps: init-submodules
	@echo "Ensuring googleapis are available..."
	@if [ ! -d "$(GOOGLEAPIS_PATH)" ]; then \
		echo "Cloning googleapis..."; \
		git clone --depth=1 https://github.com/googleapis/googleapis.git $(GOOGLEAPIS_PATH); \
	else \
		echo "googleapis already present, skipping update."; \
	fi

## proto-gen: Generate Go code from all .proto files using protoc
proto-gen:
	@echo "Generating Protobuf files via protoc..."
	mkdir -p $(GEN_OUT)
	protoc \
		--proto_path=$(PROTO_ROOT) \
		--proto_path=$(GOOGLEAPIS_PATH) \
		--go_out=$(GEN_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_OUT) --go-grpc_opt=paths=source_relative \
		$(shell find $(PROTO_ROOT) -name '*.proto' -type f)

journal-gen-proto: proto-gen
sso-gen-proto: proto-gen

journal-gen-sqlc:
	@echo "Generating SQLC for Journal Service..."
	cd $(JOURNAL_DIR) && sqlc generate

sso-gen-sqlc:
	@echo "Generating SQLC for SSO Service..."
	cd $(SSO_DIR) && sqlc generate

## proto-descriptor: Generate combined protobuf descriptor for Envoy
proto-descriptor:
	@echo "Generating protobuf descriptor..."
	protoc \
		--proto_path=$(PROTO_ROOT) \
		--proto_path=$(GOOGLEAPIS_PATH) \
		--descriptor_set_out=combined_descriptor.pb \
		--include_imports \
		$(shell find $(PROTO_ROOT) -name '*.proto' -type f)

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
	goose -table goose_sso_version -dir internal/sso_service/sql/schema postgres $(DB_URL) up
	goose -table goose_journal_version -dir internal/journal_service/sql/schema postgres $(DB_URL) up

migrate-down:
	goose -table goose_sso_version -dir internal/sso_service/sql/schema postgres $(DB_URL) down
	goose -table goose_journal_version -dir internal/journal_service/sql/schema postgres $(DB_URL) down

## help: Show available commands
help:
	@echo "Available commands:"
	@echo "  make gen          - Generate Proto and SQLC code for all services"
	@echo "  make build        - Compile binaries for all services"
	@echo "  make migrate-up   - Run database migrations for all services"
	@echo "  make migrate-down - Rollback database migrations for all services"
	@echo "  make rebuild      - Clean generated files, tidy modules, and rebuild"
