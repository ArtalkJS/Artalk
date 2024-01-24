PKG_NAME    := github.com/ArtalkJS/Artalk
BIN_NAME	:= ./bin/artalk
VERSION     ?= $(shell git describe --tags --abbrev=0)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD)

HAS_RICHGO  := $(shell which richgo)
GOTEST      ?= $(if $(HAS_RICHGO), richgo test, go test)
ARGS        ?= server

all: install build

install:
	go mod tidy

build:
	go build \
    	-ldflags "-s -w -X $(PKG_NAME)/internal/config.Version=$(VERSION) \
        -X $(PKG_NAME)/internal/config.CommitHash=$(COMMIT_HASH)" \
        -o $(BIN_NAME) \
    	$(PKG_NAME)

build-frontend:
	./scripts/build-frontend.sh

run: all
	$(BIN_NAME) $(ARGS)

build-debug:
	@echo "Building Artalk $(VERSION) for debugging..."
	@go build \
		-ldflags "-X $(PKG_NAME)/internal/config.Version=$(VERSION) \
		  -X $(PKG_NAME)/internal/config.CommitHash=$(COMMIT_HASH)" \
		-gcflags "all=-N -l" \
		-o $(BIN_NAME) \
		$(PKG_NAME)

dev: build-debug
	$(BIN_NAME) $(ARGS)

test:
	$(GOTEST) -timeout 20m ./internal/...

test-coverage:
	$(GOTEST) -cover ./...

test-coverage-html:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

update-i18n:
	go generate ./internal/i18n

update-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g server/server.go --output ./docs/swagger --requiredByDefault

docker-build:
	./scripts/docker-build.sh

docker-push:
	./scripts/docker-build.sh --push

test-frontend-e2e:
	./scripts/frontend-e2e-test.sh $(if $(REPORT), --show-report)

.PHONY: all install build build-frontend build-debug \
	dev test test-coverage test-coverage-html update-i18n \
	docker-build docker-push test-frontend-e2e;
