ifndef VERSION
	VERSION := $(shell git describe --tags --abbrev=0)
endif

COMMIT_HASH :=$(shell git rev-parse --short HEAD)
DEV_VERSION := dev-${COMMIT_HASH}

all: install update build

install:
	go mod tidy
	go install github.com/markbates/pkger/cmd/pkger

update:
	pkger -include /frontend -o http

build: update
	go build -ldflags "-X github.com/ArtalkJS/Artalk-API-Go.Version=${VERSION}" -o bin/artalk-go github.com/ArtalkJS/Artalk-API-Go

run: update build
	./bin/artalk-go

test: update
	go test -cover github.com/ArtalkJS/Artalk-API-Go/...
