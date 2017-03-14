package main

import (
	_ "fmt"
	"io"
	"log"

	pc "github.com/talbright/keds/client"
	pb "github.com/talbright/keds/gen/proto"
	ut "github.com/talbright/keds/utils/token"

	"golang.org/x/net/context"
)

const (
	address = "localhost:50051"
)

func main() {
	descriptor := &pb.PluginDescriptor{
		Name:        "example",
		Usage:       "<TODO>",
		EventFilter: "*",
		Version:     "1",
		RootCommand: "example",
	}
	client := pc.NewPluginClient(descriptor)
	if err := client.Connect(address); err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer client.Close()
	examplePlugin := NewExamplePlugin(descriptor, client)
	examplePlugin.Run()
}

type ExamplePlugin struct {
	descriptor *pb.PluginDescriptor
	client     *pc.PluginClient
}

func NewExamplePlugin(descriptor *pb.PluginDescriptor, client *pc.PluginClient) *ExamplePlugin {
	return &ExamplePlugin{
		descriptor: descriptor,
		client:     client,
	}
}

func (p *ExamplePlugin) Run() (err error) {
	ctx := ut.AddTokenToContext(context.Background(), p.client.Token)
	stream, err := p.client.EventBus(ctx)
	if err != nil {
		log.Fatalf("event bus error: %v", err)
	}
	p.client.Printf("client connected to bus")
	waitc := make(chan struct{})
	//receive
	go func() {
		for {
			var in *pb.PluginEvent
			if in, err = stream.Recv(); err == nil {
				log.Printf("received event: %v", in)
			} else if err == io.EOF {
				log.Printf("end of stream")
				close(waitc)
				return
			} else {
				log.Printf("stream recv error: %v", err)
				close(waitc)
				return
			}
		}
	}()
	<-waitc
	return
}
