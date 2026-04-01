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
	@rm -rf /tmp/vault-raft-test && mkdir -p /tmp/vault-raft-test
	@nohup vault server -config=util/fixtures/vault-test-config.hcl &
	@sleep 2
	@go test -v -run TestBootstrap ./util

shutdown:
	@killall vault

install:
	@go install .

unit:
	@go test -v ./util ./enum ./storage ./vault

accept:
	@go test -v .
