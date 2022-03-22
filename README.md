# POC gRPC API for juju

The idea of this is to showcase gRPC by creating a proxy gRPC server for the
juju client.  Documentation and client APIs for Go and Python are generated from
the protobuf specification.

## Requirements

Tested on Linux / MacOS with python3 / pip3.

## The protobuf spec

A service spec is defined in
[application.proto](./apis/juju/client/application/v1/application.proto).  The
contents of `gen/proto` is generated from this.  It contains

- a `go/` directory with
    - Generated code for Go client / server (generated with `protoc-gen-go` and
      `protoc-gen-go-grpc`);
    - Generated code for an HTTP gateway to a gRPC server (generated with
      `protoc-gen-grpc-gateway`);
- a `python/` directory with generated code for Python client (generated with
  `grpcio-tools`);
- a `doc/` directory with generated HTML documentation (with `protoc-gen-doc`).
- an `openapi/` directory with generated swagger specs

## Generating the python / Go code and HTML docs

The [buf](https://docs.buf.build/) is used to generate output from the protobuf
specs as it streamlines the process a bit.

Some protoc dependencies need to be installed before generation can be run.
```shell
make install-deps
```

Then to re-generate the code / docs:

```shell
make generate
```

## The gRPC server

This is implemented in cmd/grpc-juju.  Run it and it starts
- a gRPC server on port 8080 and
- a gRPC gateway (HTTP server) on port 8090.

It can be installed directly without downloading the repo:

```shell
go install github.com/arnodel/grpc-juju-client/cmd/grpc-juju
```

It requires a juju client to be installed somewhere on the same machine.  By
default it assumes `juju` is in the path.  You can pass in a custom path:

```shell
grpc-juju -juju-client ~/bin/juju
```

See `apis/` for supported apis.  Also the gRPC server supports reflection, so it
is discoverable (use e.g. https://github.com/fullstorydev/grpcui to play with
it)

## Example clients

The `clients/` directory contains two sample programming acting as clients to
the gRPC server.  In order to use them `grpc-juju` should be running on the same
machine.

- `clients/go-client` contains a [Go client](./clients/go-client/main.go) that
  can deploy or remove a deployment of postgresql and print the outcome.
- `clients/python-client` contains a [Python
  client](./clients/python-client/client.py) with the same functionality.  It
  should be run via the `client.sh` script to get the import paths correct.
