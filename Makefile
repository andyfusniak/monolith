OUTPUT_DIR=$(shell go env GOPATH)/bin
GITHUB_ACCOUNT=andyfusniak
GIT_COMMIT=$(shell git rev-parse --short HEAD)
PROJECT_NAME=monolith
VERSION=v0.1.0

all: monolith

monolith:
	@go build -o $(OUTPUT_DIR)/monolith -ldflags "-X 'main.version=${VERSION}' -X 'main.gitCommit=${GIT_COMMIT}'" ./cmd/monolith/main.go

dockerbuild:
	docker build \
	-t ghcr.io/${GITHUB_ACCOUNT}/${PROJECT_NAME}/monolith:latest \
	-t ghcr.io/${GITHUB_ACCOUNT}/${PROJECT_NAME}/monolith:${VERSION} \
	-f Dockerfile .

.PHONY: clean
clean:
	-@rm -r $(OUTPUT_DIR)/* 2> /dev/null || true
