package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os/exec"

	applicationv1 "github.com/arnodel/grpc-juju-client/gen/proto/go/juju/client/application/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func main() {
	var jujuClient string
	flag.StringVar(&jujuClient, "juju-client", "juju", "Path to juju client")
	flag.Parse()
	ctx := context.Background()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	log.Printf("Using juju client: %s", jujuClient)
	jujuServer := &server{
		juju: jujuClient,
	}
	applicationv1.RegisterApplicationServiceServer(grpcServer, jujuServer)
	go func() {
		log.Println("Serving gPRC on :8080")
		log.Fatal(grpcServer.Serve(lis))
	}()
	conn, err := grpc.DialContext(ctx, "0.0.0.0:8080", grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %s", err)
	}
	mux := runtime.NewServeMux()
	err = applicationv1.RegisterApplicationServiceHandler(ctx, mux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %s", err)
	}
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}
	log.Println("Serving gRPC gateway on :8090")
	log.Fatal(gwServer.ListenAndServe())
}

type server struct {
	applicationv1.UnimplementedApplicationServiceServer
	juju string
}

func (s *server) Deploy(req *applicationv1.DeployRequest, stream applicationv1.ApplicationService_DeployServer) error {
	args := []string{"deploy"}
	if req.Channel != "" {
		args = append(args, "--channel", req.Channel)
	}
	if req.Series != "" {
		args = append(args, "--series", req.Series)
	}
	if req.DryRun {
		args = append(args, "--dry-run")
	}
	if req.ArtifactName == "" {
		return status.Error(codes.InvalidArgument, "artifact name is required")
	}
	args = append(args, req.ArtifactName)
	if req.ApplicationName != "" {
		args = append(args, req.ApplicationName)
	}
	return s.runCommand(stream.Context(), stream.Send, args...)
}

func (s *server) Remove(req *applicationv1.RemoveRequest, stream applicationv1.ApplicationService_RemoveServer) error {
	args := []string{"remove-application"}
	if req.Force {
		args = append(args, "--force")
	}
	if req.NoWait {
		args = append(args, "--no-wait")
	}
	if req.ApplicationName == "" {
		return status.Error(codes.InvalidArgument, "application name is required")
	}
	args = append(args, req.ApplicationName)
	return s.runCommand(stream.Context(), stream.Send, args...)
}
func (s *server) runCommand(ctx context.Context, send func(*applicationv1.ResponseLine) error, options ...string) error {
	cmd := exec.CommandContext(ctx, s.juju, options...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	lines := make(chan *applicationv1.ResponseLine)

	runCtx, cancelRun := context.WithCancel(ctx)
	defer cancelRun()

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			select {
			case lines <- &applicationv1.ResponseLine{
				Type:    applicationv1.ResponseLineType_STDERR,
				Content: scanner.Text(),
			}:
			case <-runCtx.Done():
				return
			}
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			select {
			case lines <- &applicationv1.ResponseLine{
				Type:    applicationv1.ResponseLineType_STDOUT,
				Content: scanner.Text(),
			}:
			case <-runCtx.Done():
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case line := <-lines:
				if err := send(line); err != nil {
					// TODO: deal with error sensibly
					log.Printf("Error sending line: %s", err)
				}
			case <-runCtx.Done():
				return
			}
		}
	}()
	return cmd.Wait()
}
