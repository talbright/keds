package main

import (
	"github.com/talbright/keds/client"
	"golang.org/x/net/context"
)

type ConsoleNotifier struct {
}

func NewConsoleNotifier(client *client.Client) *ConsoleNotifier {
	return &ConsoleNotifier{}
}

func (n *ConsoleNotifier) Notify(ctx context.Context, cli *client.Client, notification *Notification) {
	if data := notification.GetData("text/plain"); data != "" {
		cli.Printf("[%s] %s", notification.Category, data)
	}
}
