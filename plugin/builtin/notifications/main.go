package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/talbright/keds/client"
	pb "github.com/talbright/keds/gen/proto"
)

type NotificationHandler struct {
}

func (e *NotificationHandler) OnCommandInvoked(cli *client.Client, event *pb.PluginEvent, cmd *cobra.Command, args []string) {
	log.Println("NotificationHandler.OnCommandInvoked")
}

func (e *NotificationHandler) OnBusEvent(cli *client.Client, event *pb.PluginEvent) {
	log.Println("NotificationHandler.OnBusEvent")
}

func (e *NotificationHandler) OnInitCommand(cli *client.Client, cmd *cobra.Command) {
	log.Println("NotificationHandler.OnInitCommand")
}

func (e *NotificationHandler) OnQuit(cli *client.Client) {
	log.Println("NotificationHandler.OnQuit")
}

func main() {
	handler := &NotificationHandler{}
	notifications, err := client.NewClient(handler)
	if err != nil {
		log.Fatalf("client error: %v", err)
	}
	notifications.Run()
}
