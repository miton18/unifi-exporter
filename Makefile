GOPATH	:= $(shell go env GOPATH)

COMMIT	:= $(shell git rev-parse HEAD)
VERSION := $(shell git describe --tags)
DATE 		:= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

build:
	go build \
	-v \
	-ldflags=" \
		-X main.version=$(VERSION) \
		-X main.commit=$(COMMIT) \
		-X main.date=$(DATE) \
	"

release_dep:
	go get -u github.com/goreleaser/goreleaser

release:
	goreleaser --skip-publish --rm-dist

release-dev:
	goreleaser --skip-publish --rm-dist --snapshot

dep:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.15.0

format:
	go fmt ./...
lint:
	./bin/golangci-lint run ./...
