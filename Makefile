BUILD_OS ?= $(shell go env GOOS)
BUILD_ARCH ?= $(shell go env GOARCH)
DOCKER_PLATFORM = $(BUILD_OS)/$(BUILD_ARCH)
PACKAGE_VERSION ?= dev

lint:
	@echo "Running linters"
	@golangci-lint run ./...

unit-test:
	@echo "Running unit tests..."
	@go test ./...

test:
	@$(MAKE) unit-test

.PHONY: build
build:
	@$(MAKE) build-platform BUILD_OS=linux BUILD_ARCH=amd64
	@$(MAKE) build-platform BUILD_OS=linux BUILD_ARCH=arm64
	@$(MAKE) build-platform BUILD_OS=darwin BUILD_ARCH=amd64
	@$(MAKE) build-platform BUILD_OS=darwin BUILD_ARCH=arm64

build-platform:
	@GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) CGO_ENABLED=0 go build -o build/httpmole-$(BUILD_OS)-$(BUILD_ARCH) cmd/httpmole/main.go

package:
	@$(MAKE) package-platform DOCKER_PLATFORM=linux/amd64
	@$(MAKE) package-platform DOCKER_PLATFORM=linux/arm64

package-platform:
	@docker build --platform $(DOCKER_PLATFORM) -t jcchavezs/httpmole:$(PACKAGE_VERSION) .

clean:
	@-rm -rf build
