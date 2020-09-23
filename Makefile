package := $(shell basename `pwd`)

.PHONY: default get codetest build setup fmt lint vet

default: fmt codetest

get:
	GOOS=windows GOARCH=amd64 go get -v ./...
	go get github.com/akavel/rsrc
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.20.0

codetest: lint vet

build: default
	mkdir -p target
	rm -f target/*
	$(shell go env GOPATH)/bin/rsrc -manifest $(package).manifest -ico $(package).ico -o $(package).syso
	GOOS=windows GOARCH=amd64 go build -v -ldflags -H=windowsgui -o target/$(package).exe

setup: build
	cp $(package).yaml target/

fmt:
	GOOS=windows GOARCH=amd64 go fmt ./...

lint:
	GOOS=windows GOARCH=amd64 $(shell go env GOPATH)/bin/golangci-lint run --fix

vet:
	GOOS=windows GOARCH=amd64 go vet -all .