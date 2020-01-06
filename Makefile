VERSION ?= dev
GIT_COMMIT ?=$(shell git rev-parse HEAD)
BUILD_DATE ?= $(shell date +%FT%T%z)
IMAGE_NAME := "jcchavezs/httpmole"

deps:
	GO111MODULES=on go get ./...

lint:
	@echo "Running linters"
	@golangci-lint run ./...

unit-test:
	@echo "Running unit tests..."
	@go test ./...

test:
	@make unit-test

build:
	go build -ldflags "-w -X main.GitCommit=${GIT_COMMIT} -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE}" -o httpmole main.go

package:
	@echo "Building image ${BIN_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) -t $(IMAGE_NAME):${VERSION} -t $(IMAGE_NAME):latest .

push:
	@echo "Pushing docker image to registry: ${VERSION} $(GIT_COMMIT)"
	docker push $(IMAGE_NAME):${VERSION}

push-latest:
	@echo "Pushing docker image to registry: latest $(GIT_COMMIT)"
	docker push $(IMAGE_NAME):latest
