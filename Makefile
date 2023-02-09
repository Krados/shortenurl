.PHONY: build
# build
build:
	go build -o ./bin/ ./...

.PHONY: install
# install
install:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/vektra/mockery/v2@latest

.PHONY: test
# test
test:
	go test -v ./... -cover

.PHONY: generate
# generate
generate:
	go generate ./...