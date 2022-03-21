install-deps:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	pip install grpcio-tools
	curl -L https://github.com/bufbuild/buf/releases/download/v1.1.0/buf-Linux-x86_64 --output buf
	chmod u+x buf

generate:
	PATH=$$PATH:$$(go env GOPATH)/bin ./buf generate