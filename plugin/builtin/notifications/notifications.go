package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/talbright/keds/client"
	pb "github.com/talbright/keds/gen/proto"
	"golang.org/x/net/context"
)

//Notifier receive events and publish them. They should create go routines
//for long running operations.
type Notifier interface {
	Notify(ctx context.Context, cli *client.Client, event *pb.PluginEvent)
}

type Notifications struct {
	notifiers []Notifier
}

func NewNotifications() *Notifications {
	return &Notifications{notifiers: make([]Notifier, 0)}
}

func (e *Notifications) AddNotifier(notifier Notifier) {
	e.notifiers = append(e.notifiers, notifier)
}

func (e *Notifications) OnCommandInvoked(cli *client.Client, event *pb.PluginEvent, cmd *cobra.Command, args []string) {
	log.Println("Notifications.OnCommandInvoked")
}

func (e *Notifications) OnBusEvent(cli *client.Client, event *pb.PluginEvent) {
	log.Println("Notifications.OnBusEvent")
	if len(e.notifiers) > 0 && e.isNotificationEvent(event) {
		e.dispatch(cli, event)
	}
}

func (e *Notifications) OnInitRootCommand(cli *client.Client, cmd *cobra.Command) {
	log.Println("Notifications.OnInitCommand")
}

func (e *Notifications) OnConnected(cli *client.Client) {
	log.Println("Notifications.OnConnected")
}

func (e *Notifications) OnRegistered(cli *client.Client) {
	log.Println("Notifications.OnRegistered")
	e.LoadNotifiers(cli)
}

func (e *Notifications) OnQuit(cli *client.Client) {
	log.Println("Notifications.OnQuit")
}

func (e *Notifications) LoadNotifiers(cli *client.Client) {
}

func (e *Notifications) isNotificationEvent(event *pb.PluginEvent) bool {
	return true
}

func (e *Notifications) dispatch(cli *client.Client, event *pb.PluginEvent) {
	for _, n := range e.notifiers {
		ctx := context.Background()
		n.Notify(ctx, cli, event)
	}
}
