SHELL=/bin/bash

dep_release_tag := v0.4.1
dep_version := $(shell dep version | grep -v -F "go version" | grep -F "version" | awk '{print $$3}')

all: unit_test build system_test

unit_test:
	go test ./...

build_all: build_linux build_darwin build_windows build_freebsd

build_linux:
	GOOS=linux GOARCH=amd64 go build -o ./artefacts/blondie_linux_amd64 cmd/blondie/main.go

build_darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./artefacts/blondie_darwin_amd64 cmd/blondie/main.go

build_windows:
	GOOS=windows GOARCH=amd64 go build -o ./artefacts/blondie_windows_amd64 cmd/blondie/main.go

build_freebsd:
	GOOS=freebsd GOARCH=amd64 go build -o ./artefacts/blondie_freebsd_amd64 cmd/blondie/main.go

build:
	go build -o ./artefacts/blondie cmd/blondie/main.go

tests:
	go test ./...

system_test: build
	go test -tags cli_tests ./cli_tests