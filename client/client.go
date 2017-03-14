package client

import (
	"fmt"

	_ "github.com/davecgh/go-spew/spew"
	pb "github.com/talbright/keds/gen/proto"
	ut "github.com/talbright/keds/utils/token"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

/*
PluginClient is used by plugins to interact with the core runtime over gRPC
*/
type PluginClient struct {
	Conn               *grpc.ClientConn
	PluginDescriptor   *pb.PluginDescriptor
	Token              string
	consoleWriteStream pb.KedsService_ConsoleWriterClient
	rootContext        context.Context
	pb.KedsServiceClient
}

func NewPluginClient(descriptor *pb.PluginDescriptor) *PluginClient {
	return &PluginClient{PluginDescriptor: descriptor, rootContext: context.Background()}
}

func (c *PluginClient) Connect(address string, opts ...grpc.DialOption) (err error) {
	if len(opts) < 1 {
		opts = append(opts, grpc.WithInsecure())
	}
	if c.Conn, err = grpc.Dial(address, opts...); err == nil {
		c.KedsServiceClient = pb.NewKedsServiceClient(c.Conn)
		err = c.register()
	}
	return
}

func (c *PluginClient) Printf(format string, args ...interface{}) (err error) {
	if c.consoleWriteStream == nil {
		c.consoleWriteStream, err = c.ConsoleWriter(c.ContextWithToken())
	}
	if err != nil {
		return
	}
	m := &pb.ConsoleWriteRequest{
		Data: fmt.Sprintf(format, args...),
	}
	err = c.consoleWriteStream.Send(m)
	return
}

func (c *PluginClient) ContextWithToken() context.Context {
	return ut.AddTokenToContext(c.rootContext, c.Token)
}

func (c *PluginClient) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
	if c.consoleWriteStream != nil {
		c.consoleWriteStream.CloseAndRecv()
	}
}

func (c *PluginClient) register() (err error) {
	registerPluginRequest := &pb.RegisterPluginRequest{
		PluginDescriptor: c.PluginDescriptor,
	}
	var header metadata.MD
	if _, err := c.RegisterPlugin(
		context.Background(),
		registerPluginRequest,
		grpc.Header(&header)); err == nil {
		c.Token = ut.GetTokenFromMetadata(header)
	}
	return
}
