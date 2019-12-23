VERSION ?= dev
GIT_COMMIT ?=$(shell git rev-parse HEAD)
BUILD_DATE ?= $(shell date +%FT%T%z)
IMAGE_NAME := "httpmole"

build:
	go build -ldflags "-w -X main.GitCommit=${GIT_COMMIT} -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE}" -o httpmole main.go

package:
	@echo "Building image ${BIN_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) -t $(IMAGE_NAME):${VERSION} -t $(IMAGE_NAME):latest .
	