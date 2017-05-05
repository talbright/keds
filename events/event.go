package events

import (
	pb "github.com/talbright/keds/gen/proto"
	"github.com/talbright/keds/plugin"
)

type IEvent interface{}

type Event struct {
	*pb.PluginEvent
}

func CobraCommandInvokedEvent(plug plugin.IPlugin, args []string) (event *pb.PluginEvent) {
	event = &pb.PluginEvent{
		Source: "keds.cli",
		Target: plug.GetName(),
		Args:   args,
	}
	return
}
