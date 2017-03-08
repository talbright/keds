package main

import (
	"log"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	. "github.com/talbright/keds/gen/minipod"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	gRPCPort      = ":50051"
	gRPCEndPoint  = "localhost:50051"
	proxyPort     = ":8080"
	debugEndPoint = "localhost:8081"
)

var idAutoIncrement uint32 = 0

// server implements minipod.MinipodServiceServer
type server struct{}

func (s *server) ProvisionMinipod(ctx context.Context, mp *ProvisionMinipodRequest) (*ProvisionMinipodResponse, error) {
	trace.FromContext(ctx)
	atomic.AddUint32(&idAutoIncrement, 1)
	return &ProvisionMinipodResponse{Id: idAutoIncrement}, nil
}

func startRpcServer() {
	log.Printf("starting rpc server")
	lis, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterMinipodServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start RPC server: %v", err)
	}
	log.Printf("stopped rpc server")
}

//curl -X POST -k http://localhost:8080/v1/minipod -d '{"name": "my name"}'
func startReverseProxy() {
	log.Printf("starting reverse proxy")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := RegisterMinipodServiceHandlerFromEndpoint(ctx, mux, gRPCEndPoint, opts); err != nil {
		log.Fatalf("failed to start reverse proxy: %v", err)
	}

	if err := http.ListenAndServe(proxyPort, mux); err != nil {
		log.Fatalf("failed to start reverse proxy: %v", err)
	}
	log.Printf("stopped reverse proxy")
}

/*
	Package golang.org/x/net/trace exports two http handlers for tracing:

	/debug/requests
	/debug/events

	ex: http://localhost:8081/debug/requests
*/
func startDebugServer() {
	if err := http.ListenAndServe(debugEndPoint, nil); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func main() {
	go startRpcServer()
	go startReverseProxy()
	go startDebugServer()
	select {}
}
