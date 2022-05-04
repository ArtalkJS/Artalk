PACKAGE_NAME := github.com/ArtalkJS/ArtalkGo
VERSION      ?= $(shell git describe --tags --abbrev=0)
COMMIT_HASH  := $(shell git rev-parse --short HEAD)
DEV_VERSION  := dev-${COMMIT_HASH}
GO_VERSION   ?= 1.18.1

HAS_RICHGO   := $(shell which richgo)
GOTEST       ?= $(if $(HAS_RICHGO), richgo test, go test)

all: install build

install:
	go mod tidy
	go install github.com/markbates/pkger/cmd/pkger

build: build-frontend update
	go build \
    	-ldflags "-s -w -X github.com/ArtalkJS/ArtalkGo/lib.Version=${VERSION} \
        -X github.com/ArtalkJS/ArtalkGo/lib.CommitHash=${COMMIT_HASH}" \
        -o bin/artalk-go \
    	github.com/ArtalkJS/ArtalkGo

build-frontend:
	./scripts/build-frontend.sh

update:
	pkger -include /frontend -include /email-tpl -include /lib/captcha/pages -include /artalk-go.example.yml -o pkged

run: all
	./bin/artalk-go server $(ARGS)

dev:
	@if [ ! -f "pkged/pkged.go" ]; then \
		make install; \
		make update; \
	fi
	@go build \
    	-ldflags "-s -w -X github.com/ArtalkJS/ArtalkGo/lib.Version=${VERSION} \
        -X github.com/ArtalkJS/ArtalkGo/lib.CommitHash=${COMMIT_HASH}" \
        -o bin/artalk-go \
    	github.com/ArtalkJS/ArtalkGo
	./bin/artalk-go server $(ARGS)

test:
	$(GOTEST) -timeout 20m ./model/...

test-coverage:
	$(GOTEST) -cover ./...

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

.PHONY: all install build build-frontend \
	update run dev test test-coverage \
	docker-build docker-push \
	release-dry-run release;
