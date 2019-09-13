package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/talbright/keds/client"
	pb "github.com/talbright/keds/gen/proto"
)

type Example struct {
	flagQuit bool
}

func (e *Example) OnCommandInvoked(cli *client.Client, event *pb.PluginEvent, cmd *cobra.Command, args []string) {
	log.Println("Plugin.Example.OnCommandInvoked")
	if e.flagQuit {
		cli.Quit(0)
	}
}

func (e *Example) OnBusEvent(cli *client.Client, event *pb.PluginEvent) {
	log.Println("Plugin.Example.OnBusEvent")
}

func (e *Example) OnInitRootCommand(cli *client.Client, cmd *cobra.Command) {
	log.Println("Plugin.Example.OnInitCommand")
	cmd.Flags().BoolVarP(&e.flagQuit, "quit", "q", false, "immediately exit")
}

func (e *Example) OnConnected(cli *client.Client) {
	log.Println("Plugin.Example.OnConnected")
}

func (e *Example) OnRegistered(cli *client.Client) {
	log.Println("Plugin.Example.OnRegistered")
}

func (e *Example) OnQuit(cli *client.Client) {
	log.Println("Plugin.Example.OnQuit")
}

func main() {
	handler := &Example{}
	example, err := client.NewClient(handler)
	if err != nil {
		log.Fatalf("client error: %v", err)
	}
	example.Run()
}
