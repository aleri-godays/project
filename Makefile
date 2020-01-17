VERSION := `git describe --tags`
SOURCES ?= $(shell find . -name "*.go" -type f)
BINARY_NAME = project
NOW = `date +"%Y-%m-%d_%H-%M-%S"`
MAIN_GO_PATH=cmd/project/main.go

all: vet lint test build build-linux

.PHONY: build
build:
	CGO_ENABLED=0 go build -i -v -o release/${BINARY_NAME} -ldflags="-X main.version=${VERSION}" ${MAIN_GO_PATH}

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i -v -o release/linux-amd64/${BINARY_NAME} -ldflags="-X main.version=${VERSION}" ${MAIN_GO_PATH}

.PHONY: test
test:
	go test -race ./...

.PHONY: cover
cover:
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: lint
lint:
	@for file in ${SOURCES} ;  do \
		golint $$file ; \
	done

.PHONY: gen
gen:
	rm -f .mock/*
	go generate ./...

.PHONY: clean
clean:
	rm -rf release/*
	rm -f cover.out
	go clean -testcache
