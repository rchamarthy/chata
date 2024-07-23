export GO111MODULE=on
GO_SRC=$(shell find . -path ./.build -prune -false -o -name \*.go)

.PHONY: all
all: lint test

test: $(GO_SRC)
	go test -v -race -cover -coverprofile=coverage.txt -covermode=atomic ./...

lint: ./.golangcilint.yaml
	./bin/golangci-lint --version || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.59.1
	./bin/golangci-lint --config ./.golangcilint.yaml run ./...

.PHONY: clean
clean:
	rm -rf bin
	rm coverage.txt