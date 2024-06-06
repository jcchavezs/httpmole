lint:
	@echo "Running linters"
	@golangci-lint run ./...

unit-test:
	@echo "Running unit tests..."
	@go test ./...

test:
	@make unit-test

.PHONY: build
build:
	CGO_ENABLED=0 go build -o build/httpmole cmd/httpmole/main.go

package:
	@docker build -t jcchavezs/httpmole:dev .

clean:
	@-rm -rf build
