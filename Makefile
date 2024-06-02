PKG_NAME    := github.com/ArtalkJS/Artalk
BIN_NAME	:= ./bin/artalk
VERSION     ?= $(shell git describe --tags --abbrev=0 --match 'v*')
COMMIT_HASH ?= $(shell git rev-parse --short HEAD)

HAS_RICHGO  := $(shell which richgo)
GOTEST      ?= $(if $(HAS_RICHGO), richgo test, go test)
ARGS        ?= server

all: install build

install:
	go mod tidy

run: all
	$(BIN_NAME) $(ARGS)

build:
	go build \
    	-ldflags "-s -w -X $(PKG_NAME)/internal/config.Version=$(VERSION) \
        -X $(PKG_NAME)/internal/config.CommitHash=$(COMMIT_HASH)" \
        -o $(BIN_NAME) \
    	$(PKG_NAME)

build-frontend:
	./scripts/build-frontend.sh

build-debug:
	@echo "Building Artalk $(VERSION) for debugging..."
	@go build \
		-ldflags "-X $(PKG_NAME)/internal/config.Version=$(VERSION) \
		  -X $(PKG_NAME)/internal/config.CommitHash=$(COMMIT_HASH)" \
		-gcflags "all=-N -l" \
		-o $(BIN_NAME) \
		$(PKG_NAME)

dev: build-debug
	ATK_TRUSTED_DOMAINS="http://localhost:5173 http://localhost:23367" \
	$(BIN_NAME) $(ARGS)

test:
	$(GOTEST) -timeout 20m $(or $(TEST_PATHS), ./...)

test-coverage:
	$(GOTEST) -cover $(or $(TEST_PATHS), ./...)

test-coverage-html:
	$(GOTEST) -v -coverprofile=coverage.out $(or $(TEST_PATHS), ./...)
	go tool cover -html=coverage.out

test-frontend-e2e:
	./scripts/frontend-e2e-test.sh $(if $(REPORT), --show-report)

update-i18n:
	go generate ./internal/i18n

update-conf:
	go generate ./internal/config

update-conf-docs:
	go run ./internal/config/meta/gen --format markdown --locale zh-cn -o ./docs/docs/guide/env.md

update-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g server/server.go --output ./docs/swagger --requiredByDefault
	pnpm -r swagger:build-http-client

docker-build:
	./scripts/docker-build.sh

docker-push:
	./scripts/docker-build.sh --push

.PHONY: all install run build build-frontend build-debug dev \
	test test-coverage test-coverage-html test-frontend-e2e \
	update-i18n update-conf update-conf-docs update-swagger \
	docker-build docker-push;
