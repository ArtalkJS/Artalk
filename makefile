ifndef VERSION
	VERSION := $(shell git describe --tags --abbrev=0)
endif

COMMIT_HASH := $(shell git rev-parse --short HEAD)
DEV_VERSION := dev-${COMMIT_HASH}

all: install update build

install:
	go mod tidy
	go install github.com/markbates/pkger/cmd/pkger

update:
	pkger -include /frontend -include /email-tpl -o pkged

build: update
	go build \
    	-ldflags "-X github.com/ArtalkJS/ArtalkGo/lib.Version=${VERSION} \
        -X github.com/ArtalkJS/ArtalkGo/lib.LastCommit=${COMMIT_SHA}" \
        -o bin/artalk-go \
    	github.com/ArtalkJS/ArtalkGo

run: update build
	./bin/artalk-go

test: update
	go test -cover github.com/ArtalkJS/ArtalkGo/...

docker-build:
	./scripts/docker-build.sh

docker-push:
	./scripts/docker-build.sh --push
