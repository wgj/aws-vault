BIN=aws-vault
OS=$(shell uname -s)
ARCH=$(shell uname -m)
GOVERSION=$(shell go version)
GOBIN=$(shell go env GOBIN)
VERSION=$(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION) -s -w
DEBUG_FLAGS=-X main.Version=$(VERSION)
CERT="Developer ID Application: 99designs Inc (NRM9HVJ62Z)"
SRC=$(shell find . -name '*.go')

test:
	go test -v $(shell go list ./... | grep -v /vendor/)

build:
	go build -o aws-vault -ldflags="$(FLAGS)" .

install:
	go install -ldflags="$(FLAGS)" .

sign:
	codesign -s $(CERT) ./aws-vault

$(BIN)-linux-amd64: $(SRC)
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags="$(FLAGS)" .

$(BIN)-darwin-amd64: $(SRC)
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags="$(FLAGS)" .

$(BIN)-darwin-amd64-debug: $(SRC)
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags="$(DEBUG_FLAGS)" .

$(BIN)-windows-386.exe: $(SRC)
	GOOS=windows GOARCH=386 go build -o $@ -ldflags="$(FLAGS)" .

release: $(BIN)-linux-amd64 $(BIN)-darwin-amd64 $(BIN)-windows-386.exe
	codesign -s $(CERT) $(BIN)-darwin-amd64

clean:
	rm -f $(BIN)-*-*
