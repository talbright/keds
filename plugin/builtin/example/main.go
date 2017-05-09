package main

import (
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"
	pc "github.com/talbright/keds/client"
	"github.com/talbright/keds/events"
	pb "github.com/talbright/keds/gen/proto"
	ut "github.com/talbright/keds/utils/token"
	"golang.org/x/net/context"
)

const (
	address = "localhost:50051"
)

var rootCmd *cobra.Command
var flagQuit bool

func init() {
	rootCmd = &cobra.Command{
		Use:   "example",
		Short: "An example plugin for the keds framework.",
		Long:  "See http://github.com/talbright/keds/README.md",
	}
	rootCmd.Flags().BoolVarP(&flagQuit, "quit", "q", false, "immediately exit")
}

func main() {
	descriptor := &pb.PluginDescriptor{
		Name:             "example",
		Usage:            "example",
		EventFilter:      "*",
		Version:          "1",
		RootCommand:      "example",
		ShortDescription: "This is an example plugin for the keds framework",
		LongDescription: `This is an exmaple plugin for the keds framework. 
Keds is a general purpose and opinionated CLI plugin framework. Plugins
communicate with the server over gRPC`,
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

func (p *ExamplePlugin) Quit(stream pb.KedsService_EventBusClient, code int) {
	quitEvent := events.CreateEventServerQuit(nil, code).Proto()
	if err := stream.Send(quitEvent); err != nil {
		log.Printf("failed to send event: %s", err)
	}
}

func (p *ExamplePlugin) Run() (err error) {
	ctx := ut.AddTokenToContext(context.Background(), p.client.Token)
	stream, err := p.client.EventBus(ctx)
	if err != nil {
		log.Fatalf("event bus error: %v", err)
	}
	//writing to stdout/stderr works as well
	p.client.Printf("example plugin connected to console")
	waitc := make(chan struct{})
	//event loop
	go func() {
		for {
			if in, err := stream.Recv(); err == nil {
				log.Printf("plugin received event: %v", in)
				if in.GetName() == "keds.command_invoked" {
					rootCmd.SetArgs(in.GetArgs())
					rootCmd.Run = func(cmd *cobra.Command, args []string) {
						log.Printf("Cobra.Run")
						if flagQuit {
							p.Quit(stream, 0)
						}
					}
					if err := rootCmd.Execute(); err != nil {
						fmt.Println("error for rootCmd")
						p.Quit(stream, 2)
					}
				}
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
