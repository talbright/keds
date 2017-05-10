package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/talbright/keds/client"
	pb "github.com/talbright/keds/gen/proto"
)

type ExampleHandler struct {
	flagQuit bool
}

func (e *ExampleHandler) OnCommandInvoked(cli *client.Client, event *pb.PluginEvent, cmd *cobra.Command, args []string) {
	log.Println("ExampleHandler.OnCommandInvoked")
	if e.flagQuit {
		cli.Quit(0)
	}
}

func (e *ExampleHandler) OnBusEvent(cli *client.Client, event *pb.PluginEvent) {
	log.Println("ExampleHandler.OnBusEvent")
}

func (e *ExampleHandler) OnInitCommand(cli *client.Client, cmd *cobra.Command) {
	log.Println("ExampleHandler.OnInitCommand")
	cmd.Flags().BoolVarP(&e.flagQuit, "quit", "q", false, "immediately exit")
}

func (e *ExampleHandler) OnQuit(cli *client.Client) {
	log.Println("ExampleHandler.OnQuit")
}

func main() {
	handler := &ExampleHandler{}
	example, err := client.NewClient(handler)
	if err != nil {
		log.Fatalf("client error: %v", err)
	}
	example.Run()
}
