package main

import (
	"log"

	. "github.com/talbright/keds/gen/minipod"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {

	//connection to server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//invoke client
	c := NewMinipodServiceClient(conn)
	r, err := c.Create(context.Background(), &Minipod{Name: "myname"})
	if err != nil {
		log.Fatalf("request error: %v", err)
	}
	log.Printf("response: %v", r)
}
