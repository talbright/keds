package main

import (
	"github.com/talbright/keds/client"
	"golang.org/x/net/context"
)

//ConsoleNotifier sends any text/plain messages it receives to the console
//using the RPC client Printf method.
type ConsoleNotifier struct{}

//NewConsoleNotifier creates a console notifier that implements the Notifier
//interface.
func NewConsoleNotifier(client *client.Client) *ConsoleNotifier {
	return &ConsoleNotifier{}
}

//Notify implementation of Notifier interface
func (n *ConsoleNotifier) Notify(ctx context.Context, cli *client.Client, notification *Notification) {
	if data := notification.GetData("text/plain"); data != "" {
		cli.Printf("[%s] %s", notification.Category, data)
	}
}
