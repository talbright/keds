package main

import (
	"github.com/talbright/keds/plugin"
	. "github.com/talbright/keds/server"
)

func main() {

	pluginRegistry := plugin.NewPluginRegistry()
	eventBusAdapter := plugin.NewEventBus()
	gRPC := &KedsRPCServer{
		Plugins:         pluginRegistry,
		EventBusAdapter: eventBusAdapter,
	}
	go gRPC.Start()
	select {}

}
