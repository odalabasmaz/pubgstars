# note: call scripts from /scripts

.PHONY: setup build test

setup:
	git config core.hooksPath ../.githooks

build:
	go build ./...

test:
	go test ./...