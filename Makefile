.PHONY: build

fmt:
	@go fmt ./...

tidy:
	@go mod tidy

build: tidy
	@go build -o vault-raft-backup main.go

release: tidy
	@go build -ldflags="-s -w" -o vault-raft-backup main.go

bootstrap:
	@rm -f nohup.out
	# using cli for this avoids importing the entire vault/command package
	@nohup vault server -dev -dev-root-token-id="abcdefghijklmnopqrstuvwxyz09" &
	@go test -v -run TestBootstrap ./util

shutdown:
	@killall vault

install:
	@go install .

unit:
	@go test -v ./storage ./util ./vault

accept:
	@go test -v .
