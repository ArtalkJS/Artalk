PACKAGE_NAME := github.com/ArtalkJS/Artalk
BIN_NAME	 := ./bin/artalk
VERSION      ?= $(shell git describe --tags --abbrev=0)
COMMIT_HASH  := $(shell git rev-parse --short HEAD)
DEV_VERSION  := dev-${COMMIT_HASH}
GO_VERSION   ?= 1.19.4

HAS_RICHGO   := $(shell which richgo)
GOTEST       ?= $(if $(HAS_RICHGO), richgo test, go test)

all: install build

install:
	go mod tidy

build: build-frontend
	go build \
    	-ldflags "-s -w -X github.com/ArtalkJS/Artalk/internal/config.Version=${VERSION} \
        -X github.com/ArtalkJS/Artalk/internal/config.CommitHash=${COMMIT_HASH}" \
        -o $(BIN_NAME) \
    	github.com/ArtalkJS/Artalk

build-frontend:
	./scripts/build-frontend.sh

run: all
	$(BIN_NAME) server $(ARGS)

debug-build:
	@if [ ! -f "pkged/pkged.go" ]; then \
		make install; \
	fi
	@echo "Building Artalk ${VERSION} for debugging..."
	@go build \
		-ldflags " \
			-X github.com/ArtalkJS/Artalk/internal/config.Version=${VERSION} \
			-X github.com/ArtalkJS/Artalk/internal/config.CommitHash=${COMMIT_HASH}" \
		-gcflags "all=-N -l" \
		-o $(BIN_NAME) \
		github.com/ArtalkJS/Artalk

dev: debug-build
	$(BIN_NAME) server $(ARGS)

test:
	$(GOTEST) -timeout 20m ./internal/...

test-coverage:
	$(GOTEST) -cover ./...

update-i18n:
	go generate ./internal/i18n

docker-build:
	./scripts/docker-build.sh

docker-push:
	./scripts/docker-build.sh --push

release-dry-run:
	@docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:v${GO_VERSION} \
		--rm-dist --skip-validate --skip-publish


# https://hub.docker.com/r/troian/golang-cross
# https://github.com/troian/golang-cross
# https://goreleaser.com/cmd/goreleaser_release/
# --skip-validate 参数跳过 git checks (由于 pkger 和 .release-env 文件生成)
release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:v${GO_VERSION} \
		release --rm-dist --skip-validate

.PHONY: all install build debug-build build-frontend \
	run dev test test-coverage \
	docker-build docker-push \
	release-dry-run release;
