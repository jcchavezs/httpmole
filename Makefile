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
	go build -o httpmole cmd/httpmole/main.go

clean:
	@-rm httpmole
	@-rm -rf dist

release:
	@echo "Make sure you are logged in dockerhub"
	GITHUB_TOKEN=$(GITHUB_TOKEN) goreleaser release --rm-dist

release.dryrun:
	goreleaser release --skip-publish --snapshot --rm-dist
