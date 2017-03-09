package client

import (
	"fmt"

	pb "github.com/talbright/keds/gen/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

/*
PluginClient is used by plugins to interact with the core runtime over gRPC
*/
type PluginClient struct {
	Conn               *grpc.ClientConn
	PluginDescriptor   *pb.PluginDescriptor
	Signature          string
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
		//TODO is this the correct way to use context?
		//see: https://godoc.org/golang.org/x/net/context#Context
		// var ctx context.Context
		// ctx, c.consoleWriteStreamCancelFunc = context.WithCancel(c.rootContext)
		// c.consoleWriteStream, err = c.Client.ConsoleWriter(ctx)
		c.consoleWriteStream, err = c.ConsoleWriter(context.Background())
	}
	if err != nil {
		return
	}
	m := &pb.ConsoleWriteRequest{
		Signature: c.Signature,
		Data:      fmt.Sprintf(format, args...),
	}
	err = c.consoleWriteStream.Send(m)
	return
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
	if results, err := c.RegisterPlugin(context.Background(), registerPluginRequest); err == nil {
		c.Signature = results.GetSignature()
	}
	return
}
