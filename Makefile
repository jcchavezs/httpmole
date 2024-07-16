GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

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
	@$(MAKE) build-platform GOOS=linux GOARCH=amd64
	@$(MAKE) build-platform GOOS=linux GOARCH=arm64

build-platform:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o build/httpmole-$(GOOS)-$(GOARCH) cmd/httpmole/main.go

package:
	@$(MAKE) package-platform PLATFORM=linux/amd64
	@$(MAKE) package-platform PLATFORM=linux/arm64

package-platform:
	@docker build --platform $(PLATFORM) -t jcchavezs/httpmole:dev .

clean:
	@-rm -rf build
