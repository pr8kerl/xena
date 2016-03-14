GOROOT := /usr/local/go
GOPATH := $(shell pwd)
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)
DEPS   := github.com/mitchellh/cli github.com/aws/aws-sdk-go/aws  github.com/aws/aws-sdk-go/service/ec2 github.com/aws/aws-sdk-go/service/iam
FILES  := main.go instances.go snapshots.go

all: xena

deps: $(DEPS)
	GOPATH=$(GOPATH) go get -u $^


xena: $(FILES)
    # always format code
		GOPATH=$(GOPATH) go fmt $^
    # vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOPATH=$(GOPATH) go build -o $@ -v $^
		touch $@

linux64: $(FILES)
    # always format code
		GOPATH=$(GOPATH) go fmt $^
    # vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOOS=linux GOARCH=amd64 GOPATH=$(GOPATH) go build -o xena-linux-amd64.bin -v $^
		touch xena-linux-amd64.bin

win64: $(FILES)
    # always format code
		GOPATH=$(GOPATH) go fmt $^
    # vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOOS=windows GOARCH=amd64 GOPATH=$(GOPATH) go build -o xena-win-amd64.exe -v $^
		touch xena-win-amd64.exe

.PHONY: $(DEPS) clean

clean:
		rm -f xena xena-win-amd64.exe xena-linux-amd64.bin

