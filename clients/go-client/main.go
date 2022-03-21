package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	applicationv1 "github.com/arnodel/grpc-juju-client/gen/proto/go/juju/client/application/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LineReceiver interface {
	Recv() (*applicationv1.ResponseLine, error)
}

func main() {
	var rm bool
	flag.BoolVar(&rm, "remove", false, "remove instead of deploy")
	flag.Parse()

	ctx := context.Background()

	// Get a connection to the gRPC server
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error getting connection: %s", err)
	}

	// Obtain the client
	client := applicationv1.NewApplicationServiceClient(conn)

	var r LineReceiver
	if rm {
		// Make a request to remove application
		r, err = client.Remove(ctx, &applicationv1.RemoveRequest{
			ApplicationName: "postgresql",
		})
	} else {
		// Make a deployment request
		r, err = client.Deploy(ctx, &applicationv1.DeployRequest{
			ArtifactName: "postgresql",
		})
	}
	if err != nil {
		log.Fatalf("Error calling Deploy: %s", err)
	}

	// Consume the response stream
	for {
		line, err := r.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			break
		}
		fmt.Printf("%s: %s\n", applicationv1.ResponseLineType_name[int32(line.Type)], line.Content)
	}
}
