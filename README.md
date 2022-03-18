# POC gRPC API for juju

There is a command in cmd/grpc-juju.  Run it and it starts a gRPC server on port
8080 and a gRPC gateway (HTTP server) on port 8090.

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
