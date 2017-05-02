package plugin

import (
	pb "github.com/talbright/keds/gen/proto"
)

type IEvent interface{}

type Event struct {
	*pb.PluginEvent
}

func CobraCommandInvokedEvent(plugin IPlugin, args []string) (event *pb.PluginEvent) {
	event = &pb.PluginEvent{
		Source: "keds.cli",
		Target: plugin.GetName(),
		Args:   args,
	}
	return
}
