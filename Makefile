.PHONY: build

fmt:
	@go fmt ./...

tidy:
	@go mod tidy

build: tidy
	@go build -o vault-raft-backup main.go

release: tidy
	@go build -ldflags="-s -w" -o vault-raft-backup main.go

install:
	@go install .

unit:
	@go test -v ./...
