package main

import (
	"log"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	. "github.com/talbright/keds/gen/minipod"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port     = ":50051"
	endPoint = "localhost:50051"
)

var idAutoIncrement uint32 = 0

// server implements minipod.MinipodServiceServer
type server struct{}

func (s *server) ProvisionMinipod(ctx context.Context, mp *ProvisionMinipodRequest) (*ProvisionMinipodResponse, error) {
	atomic.AddUint32(&idAutoIncrement, 1)
	return &ProvisionMinipodResponse{Id: idAutoIncrement}, nil
}

func startRpcServer() {
	log.Printf("starting rpc server")
	lis, err := net.Listen("tcp", port)
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
	if err := RegisterMinipodServiceHandlerFromEndpoint(ctx, mux, endPoint, opts); err != nil {
		log.Fatalf("failed to start reverse proxy: %v", err)
	}

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to start reverse proxy: %v", err)
	}
	log.Printf("stopped reverse proxy")
}

func main() {
	go startRpcServer()
	go startReverseProxy()
	select {}
}
