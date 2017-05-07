package server

import (
	"context"
	"fmt"
	"runtime"

	"golang.org/x/net/trace"

	"github.com/spf13/cobra"
	"github.com/talbright/keds/events"
	"github.com/talbright/keds/plugin"
)

type Cobra struct {
	events  trace.EventLog
	rootCmd *cobra.Command
}

func NewCobra(rootCmd *cobra.Command) *Cobra {
	_, file, line, _ := runtime.Caller(0)
	return &Cobra{
		rootCmd: rootCmd,
		events:  trace.NewEventLog("server.Cobra", fmt.Sprintf("%s:%d", file, line)),
	}
}

func (c *Cobra) AddPlugin(ctx context.Context, plug *plugin.Plugin, bus events.IEventBus) {
	if plug.GetRootCommand() != "" {
		cmd := &cobra.Command{
			Use:                plug.GetRootCommand(),
			Short:              plug.GetShortDescription(),
			Long:               plug.GetLongDescription(),
			DisableFlagParsing: true,
			Run: func(cmd *cobra.Command, args []string) {
				c.events.Printf("cli invocation for plugin %s with args %v", plug, args)
				event := events.CreateEventCommandInvoked(plug, args).Proto()
				bus.Publish(ctx, event)
			},
		}
		c.rootCmd.AddCommand(cmd)
	}
}
