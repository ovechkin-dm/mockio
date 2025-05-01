all: gofumpt import lint

init:
	go install mvdan.cc/gofumpt@v0.6.0
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2
	go install github.com/daixiang0/gci@v0.12.3
	go install github.com/vektra/mockery/v3@v3.2.5

lint:
	golangci-lint run ./...

gofumpt:
	gofumpt -l -w .

import:
	gci write --skip-generated -s standard -s default -s "prefix(github.com/ovechkin-dm/mockio)" -s blank -s dot -s alias .

codegen:
	mockery
