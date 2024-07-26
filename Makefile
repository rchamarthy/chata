export GO111MODULE=on
TOP=$(shell git rev-parse --show-toplevel)
GO_SRC=$(shell find $(TOP) -path ./.build -prune -false -o -name \*.go)

.PHONY: all
all: lint test server client

test: $(GO_SRC)
	cd $(TOP) && go test -v -race -cover -coverprofile=coverage.txt -covermode=atomic ./...

lint: ./.golangcilint.yaml $(GO_SRC)
	cd $(TOP) && $(TOP)/bin/golangci-lint --version || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.59.1
	cd $(TOP) && $(TOP)/bin/golangci-lint --config ./.golangcilint.yaml run ./...

server: $(GO_SRC)
	cd $(TOP)/cmd/server && go build -o $(TOP)/build/chata-server.bin

client: $(GO_SRC)
	cd $(TOP)/cmd/client && go build -o $(TOP)/build/chata

run-server: server
	$(TOP)/build/chata-server.bin test.cfg

.PHONY: clean
clean:
	rm -rf bin
	rm coverage.txt
