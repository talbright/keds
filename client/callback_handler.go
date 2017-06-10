package client

import (
	"github.com/spf13/cobra"
	"github.com/talbright/keds/gen/proto"
)

//CallbackHandler implementations should process the main logic of a plugin.
type CallbackHandler interface {
	OnCommandInvoked(client *Client, event *proto.PluginEvent, cmd *cobra.Command, args []string)
	OnBusEvent(client *Client, event *proto.PluginEvent)
	OnInitRootCommand(client *Client, cmd *cobra.Command)
	OnQuit(client *Client)
	OnRegistered(client *Client)
	OnConnected(client *Client)
}
