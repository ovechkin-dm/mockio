all: gofumpt import lint

init:
	go install mvdan.cc/gofumpt@v0.6.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2
	go install github.com/daixiang0/gci@v0.12.3

lint:
	golangci-lint run ./...

gofumpt:
	gofumpt -l -w .

import:
	gci write --skip-generated -s standard -s default -s "prefix(github.com/ovechkin-dm/mockio)" -s blank -s dot -s alias .
