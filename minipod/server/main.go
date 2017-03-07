package main

import (
	"log"
	"net"

	. "github.com/talbright/keds/gen/minipod"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server implements minipod.MinipodServiceServer
type server struct{}

func (s *server) Create(ctx context.Context, mp *Minipod) (*Minipod, error) {
	return &Minipod{Id: 1, Status: "new", Name: mp.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterMinipodServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
