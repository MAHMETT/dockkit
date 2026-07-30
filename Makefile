.PHONY: build test lint clean dev install version

dev:
	air

build:
	go build -o bin/dockkit .

test:
	go test ./...

test-cover:
	go test -cover ./... -coverprofile=coverage.out

lint:
	golangci-lint run

clean:
	rm -rf bin/ dist/ coverage.out

install:
	go install .

version:
	@echo "dockkit $(shell git describe --tags --always 2>/dev/null || echo 'dev')"

tidy:
	go mod tidy

fmt:
	gofmt -s -w .
