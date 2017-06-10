package main

import (
	"github.com/talbright/keds/client"
	"github.com/talbright/keds/gen/proto"
	"golang.org/x/net/context"
)

//Notifier receive events and publish them. They should create go routines for
//long running operations. All types that wish to receive notifcations should
//implement this interface.
type Notifier interface {
	Notify(ctx context.Context, cli *client.Client, notification *Notification)
}

//Notification provides the data necessary for a Notifier to perform
//notifications, and is derived from an underlying notifcation event received
//by the plugin.
type Notification struct {
	Category string
	*Message
}

//NewNotification creates a notification consumable by any type that implements
//the Notifer interface.
func NewNotification(event *proto.PluginEvent) *Notification {
	n := &Notification{Message: NewMessage(event)}
	if val, ok := event.Data["category"]; ok {
		n.Category = val
	} else {
		n.Category = "info"
	}
	return n
}
