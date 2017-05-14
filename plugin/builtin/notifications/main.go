package main

import (
	"log"

	"github.com/talbright/keds/client"
)

func main() {
	notifications := &Notifications{}
	plugin, err := client.NewClient(notifications)
	if err != nil {
		log.Fatalf("client error: %v", err)
	}
	plugin.Run()
}
