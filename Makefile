SHELL=/bin/bash

dep_release_tag := v0.4.1
dep_version := $(shell dep version | grep -v -F "go version" | grep -F "version" | awk '{print $$3}')

all: test build

test:
	go test ./...

build:
	GOOS=linux GOARCH=amd64 go build -o ./artefacts/blondie_linux_amd64 cmd/blondie/main.go
	GOOS=darwin GOARCH=amd64 go build -o ./artefacts/blondie_darwin_amd64 cmd/blondie/main.go
	GOOS=windows GOARCH=amd64 go build -o ./artefacts/blondie_windows_amd64 cmd/blondie/main.go

