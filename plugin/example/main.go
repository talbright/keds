package main

import (
	"fmt"
	"io"
	"log"

	pc "github.com/talbright/keds/client"
	pb "github.com/talbright/keds/gen/proto"

	"golang.org/x/net/context"
)

const (
	address = "localhost:50051"
)

func main() {
	descriptor := &pb.PluginDescriptor{
		Name:        "core",
		Usage:       "tbd",
		EventFilter: "*",
		Version:     "1",
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
	stream, err := p.client.EventBus(context.Background())
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
	//send
	go func() {
		for i := 0; i < 3; i++ {
			e := &pb.PluginEvent{
				Name:   fmt.Sprintf("plugin:example:event%d", i),
				Source: "plugin:example",
			}
			log.Printf("sending event: %v", e)
			if err := stream.Send(e); err != nil {
				log.Printf("error sending event: %v", err)
			}
		}
	}()
	<-waitc
	return
}
