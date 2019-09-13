package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/talbright/keds/client"
	"github.com/talbright/keds/gen/proto"
	"golang.org/x/net/context"
)

//NotificationsHandler implements the ClientCallbackHandler, enabling it to
//act as a keds plugin. Dispatches notification events received to various
//notifiers (types that implement Notifier.)
type NotificationsHandler struct {
	notifiers []Notifier
}

func NewNotificationsHandler() *NotificationsHandler {
	return &NotificationsHandler{notifiers: make([]Notifier, 0)}
}

func (e *NotificationsHandler) AddNotifier(notifier Notifier) {
	e.notifiers = append(e.notifiers, notifier)
}

func (e *NotificationsHandler) OnCommandInvoked(cli *client.Client, event *proto.PluginEvent, cmd *cobra.Command, args []string) {
	log.Println("Notifications.OnCommandInvoked")
}

func (e *NotificationsHandler) OnBusEvent(cli *client.Client, event *proto.PluginEvent) {
	log.Println("Notifications.OnBusEvent")
	if len(e.notifiers) > 0 && e.isNotificationEvent(event) {
		e.dispatch(cli, event)
	}
}

func (e *NotificationsHandler) OnInitRootCommand(cli *client.Client, cmd *cobra.Command) {
	log.Println("Notifications.OnInitCommand")
}

func (e *NotificationsHandler) OnConnected(cli *client.Client) {
	log.Println("Notifications.OnConnected")
}

func (e *NotificationsHandler) OnRegistered(cli *client.Client) {
	log.Println("Notifications.OnRegistered")
	e.LoadNotifiers(cli)
}

func (e *NotificationsHandler) OnQuit(cli *client.Client) {
	log.Println("Notifications.OnQuit")
}

func (e *NotificationsHandler) LoadNotifiers(cli *client.Client) {
	e.AddNotifier(NewConsoleNotifier(cli))
}

func (e *NotificationsHandler) isNotificationEvent(event *proto.PluginEvent) bool {
	return event.GetName() == "notification"
}

func (e *NotificationsHandler) dispatch(cli *client.Client, event *proto.PluginEvent) {
	for _, n := range e.notifiers {
		ctx := context.Background()
		notification := NewNotification(event)
		n.Notify(ctx, cli, notification)
	}
}
