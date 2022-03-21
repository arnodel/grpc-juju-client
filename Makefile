BUF_VERSION := 1.1.0
GO_BIN := $(shell go env GOPATH)/bin

export PATH := $(PATH):$(GO_BIN)

.PHONY: install-deps install-python-deps generate

$(GO_BIN)/protoc-gen-grpc-gateway:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway

$(GO_BIN)/protoc-gen-doc:
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

$(GO_BIN)/protoc-gen-openapiv2:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

install-python-deps:
	pip install grpcio-tools google-api-python-client-stubs

bin/buf:
	curl -sSL https://github.com/bufbuild/buf/releases/download/v$(BUF_VERSION)/buf-$(shell uname -s)-$(shell uname -m) --output bin/buf
	chmod u+x bin/buf

install-deps: $(GO_BIN)/protoc-gen-grpc-gateway
install-deps: $(GO_BIN)/protoc-gen-doc
install-deps: $(GO_BIN)/protoc-gen-openapiv2
install-deps: install-python-deps bin/buf

generate:
	./bin/buf generate