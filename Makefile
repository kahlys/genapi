VERSION = $(shell git describe --tags --always)

all: build

## build: Build the binary.
build:
	go build -o genapi -ldflags "-X main.version=${VERSION}" cmd/genapi/*.go

install:
	go install ./cmd/genapi/

## lint: Run linter on source code.
lint:
	golangci-lint run --exclude-use-default=false --disable "errcheck" --enable "goconst,goimports,golint,gofmt" ./...

## clean: Remove build files.
clean:
	rm -f genapi

help: Makefile
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'