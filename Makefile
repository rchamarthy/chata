export GO111MODULE=on
TOP=$(shell git rev-parse --show-toplevel)
GO_SRC=$(shell find $(TOP) -path ./.build -prune -false -o -name \*.go)

.PHONY: all
all: lint test

test: $(GO_SRC)
	go test -v -race -cover -coverprofile=coverage.txt -covermode=atomic ./...

lint: ./.golangcilint.yaml
	$(TOP)/bin/golangci-lint --version || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.59.1
	$(TOP)/bin/golangci-lint --config ./.golangcilint.yaml run ./...

server: $(GO_SRC)
	cd ./cmd/server && go build -o $(TOP)/build/chata-server.bin

client: $(GO_SRC)
	cd ./cmd/client && go build -o $(TOP)/build/chata

run-server: server
	$(TOP)/build/chata-server.bin test.cfg

.PHONY: clean
clean:
	rm -rf bin
	rm coverage.txt