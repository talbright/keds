/*
Package client provides the primary set of types and interfaces needed to implement
a plugin compatible with keds.
*/
package client

import (
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/talbright/keds/events"
	pb "github.com/talbright/keds/gen/proto"
	"github.com/talbright/keds/server"
	"github.com/talbright/keds/utils/config"
	ut "github.com/talbright/keds/utils/token"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

/*
ClientCallbackHandler implementations should process the main logic of a plugin.
*/
type ClientCallbackHandler interface {
	OnCommandInvoked(client *Client, event *pb.PluginEvent, cmd *cobra.Command, args []string)
	OnBusEvent(client *Client, event *pb.PluginEvent)
	OnInitRootCommand(client *Client, cmd *cobra.Command)
	OnQuit(client *Client)
	OnRegistered(client *Client)
	OnConnected(client *Client)
}

/*
Client is the primary type that should be used to create plugins. Plugins
communicate with the plugin host over gRPC.
*/
type Client struct {
	CallbackHandler    ClientCallbackHandler
	Conn               *grpc.ClientConn
	PluginDescriptor   *pb.PluginDescriptor
	Token              string
	ConsoleWriteStream pb.KedsService_ConsoleWriterClient
	EventBusStream     pb.KedsService_EventBusClient
	//RootCommand is the root command for your plugins CLI. The RootCommand is only
	//initialized if your plugin descriptor has defined a root command name.
	RootCommand *cobra.Command
	//PluginConfig should define all the properties to fulfill the PluginDescriptor. Plugin
	//users should not put custom settings here, though you can store defaults.
	PluginConfig *viper.Viper
	//ServerConfig is the complete configuration for the plugin host and all other plugins.
	//Users should configure any plugin specific behavior here, not in PluginConfig.
	ServerConfig *viper.Viper
	//KedsServiceClient is the gRPC implementation based on the proto file.
	pb.KedsServiceClient
}

func NewClient(handler ClientCallbackHandler) (client *Client, err error) {
	client = &Client{CallbackHandler: handler, PluginDescriptor: &pb.PluginDescriptor{}}
	if err = client.loadPluginConfig(); err != nil {
		return
	}
	if err = client.loadServerConfig(); err != nil {
		return
	}
	client.initRootCommand()
	return
}

func (c *Client) loadServerConfig() (err error) {
	c.ServerConfig = viper.New()
	if configFilePath := config.GetConfigPath(); configFilePath != "" {
		c.ServerConfig.SetConfigFile(configFilePath)
		err = c.ServerConfig.ReadInConfig()
	}
	return
}

func (c *Client) loadPluginConfig() (err error) {
	c.PluginConfig = viper.New()
	c.PluginConfig.SetConfigFile("plugin.yaml")
	if err = c.PluginConfig.ReadInConfig(); err != nil {
		return
	}
	c.PluginDescriptor.ShortDescription = c.PluginConfig.GetString("short_description")
	c.PluginDescriptor.LongDescription = c.PluginConfig.GetString("long_description")
	c.PluginDescriptor.RootCommand = c.PluginConfig.GetString("root_command")
	c.PluginDescriptor.Name = c.PluginConfig.GetString("name")
	c.PluginDescriptor.Usage = c.PluginConfig.GetString("usage")
	c.PluginDescriptor.Version = c.PluginConfig.GetString("version")
	c.PluginDescriptor.EventFilter = c.PluginConfig.GetString("event_filter")
	return
}

func (c *Client) initRootCommand() {
	if c.PluginDescriptor.GetRootCommand() != "" {
		c.RootCommand = &cobra.Command{
			Use:   c.PluginDescriptor.GetRootCommand(),
			Short: c.PluginDescriptor.GetShortDescription(),
			Long:  c.PluginDescriptor.GetLongDescription(),
		}
		c.CallbackHandler.OnInitRootCommand(c, c.RootCommand)
	}
}

func (c *Client) Run() (err error) {
	if err = c.connect(server.EndPoint()); err != nil {
		return
	}
	c.CallbackHandler.OnConnected(c)
	if err = c.register(); err != nil {
		return
	}
	c.CallbackHandler.OnRegistered(c)
	if c.EventBusStream, err = c.EventBus(c.contextWithToken()); err != nil {
		return
	}
	c.Printf("example plugin connected to console")
	c.loop()
	return
}

func (c *Client) connect(address string, opts ...grpc.DialOption) (err error) {
	if len(opts) < 1 {
		opts = append(opts, grpc.WithInsecure())
	}
	if c.Conn, err = grpc.Dial(address, opts...); err == nil {
		c.KedsServiceClient = pb.NewKedsServiceClient(c.Conn)
	}
	return
}

func (c *Client) Printf(format string, args ...interface{}) (err error) {
	if c.ConsoleWriteStream == nil {
		if c.ConsoleWriteStream, err = c.ConsoleWriter(c.contextWithToken()); err != nil {
			return
		}
	}
	m := &pb.ConsoleWriteRequest{
		Data: fmt.Sprintf(format, args...),
	}
	err = c.ConsoleWriteStream.Send(m)
	return
}

func (c *Client) Quit(code int) {
	quitEvent := events.CreateEventServerQuit(nil, code).Proto()
	if err := c.EventBusStream.Send(quitEvent); err != nil {
		log.Printf("failed to send event: %s", err)
	}
	c.CallbackHandler.OnQuit(c)
	c.close()
}

func (c *Client) contextWithToken() context.Context {
	return ut.AddTokenToContext(context.Background(), c.Token)
}

func (c *Client) close() {
	if c.ConsoleWriteStream != nil {
		c.ConsoleWriteStream.CloseSend()
	}
	if c.EventBusStream != nil {
		c.EventBusStream.CloseSend()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
}

func (c *Client) register() (err error) {
	registerPluginRequest := &pb.RegisterPluginRequest{
		PluginDescriptor: c.PluginDescriptor,
	}
	var header metadata.MD
	if _, err = c.RegisterPlugin(
		context.Background(),
		registerPluginRequest,
		grpc.Header(&header)); err == nil {
		c.Token = ut.GetTokenFromMetadata(header)
	}
	return
}

func (c *Client) loop() (err error) {
	waitc := make(chan struct{})
	go func() {
		for {
			if in, err := c.EventBusStream.Recv(); err == nil {
				log.Printf("plugin received event: %v", in)
				c.onBusEvent(in)
			} else if err == io.EOF {
				log.Printf("end of stream")
				close(waitc)
				return
			} else {
				log.Printf("stream recv error: %v", err)
				close(waitc)
				return
			}
		}
	}()
	<-waitc
	return
}

func (c *Client) onBusEvent(in *pb.PluginEvent) {
	c.CallbackHandler.OnBusEvent(c, in)
	if in.GetName() == "keds.command_invoked" && c.RootCommand != nil {
		c.RootCommand.SetArgs(in.GetArgs())
		c.RootCommand.Run = func(cmd *cobra.Command, args []string) {
			c.CallbackHandler.OnCommandInvoked(c, in, cmd, args)
		}
		if err := c.RootCommand.Execute(); err != nil {
			log.Printf("error executing command: %s", err)
			c.Quit(1)
		}
	}
}
