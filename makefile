PACKAGE_NAME := github.com/ArtalkJS/ArtalkGo
VERSION      ?= $(shell git describe --tags --abbrev=0)
COMMIT_HASH  := $(shell git rev-parse --short HEAD)
DEV_VERSION  := dev-${COMMIT_HASH}
GO_VERSION   ?= 1.16.6

all: install build

.PHONY: install
install:
	go mod tidy
	go install github.com/markbates/pkger/cmd/pkger

.PHONY: build
build: build-frontend update
	go build \
    	-ldflags "-s -w -X github.com/ArtalkJS/ArtalkGo/lib.Version=${VERSION} \
        -X github.com/ArtalkJS/ArtalkGo/lib.CommitHash=${COMMIT_HASH}" \
        -o bin/artalk-go \
    	github.com/ArtalkJS/ArtalkGo

.PHONY: build-frontend
build-frontend:
	./scripts/build-frontend.sh

.PHONY: update
update:
	pkger -include /frontend -include /email-tpl -o pkged

.PHONY: run
run: all
	./bin/artalk-go serve $(ARGS)

.PHONY: dev
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
	./bin/artalk-go $(ARGS)

.PHONY: test
test: update
	go test -cover github.com/ArtalkJS/ArtalkGo/...

.PHONY: docker-build
docker-build:
	./scripts/docker-build.sh

.PHONY: docker-push
docker-push:
	./scripts/docker-build.sh --push

.PHONY: release-dry-run
release-dry-run:
	@docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		troian/golang-cross:v${GO_VERSION} \
		--rm-dist --skip-validate --skip-publish

.PHONY: release
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
		troian/golang-cross:v${GO_VERSION} \
		release --rm-dist --skip-validate
