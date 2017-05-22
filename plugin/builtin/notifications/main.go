package main

import (
	"log"

	"github.com/talbright/keds/client"
)

func main() {
	notificationsHandler := NewNotificationsHandler()
	plugin, err := client.NewClient(notificationsHandler)
	if err != nil {
		log.Fatalf("client error: %v", err)
	}
	plugin.Run()
}
