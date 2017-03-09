package plugin

import (
	"crypto/sha1"
	"fmt"

	pb "github.com/talbright/keds/gen/proto"
	"golang.org/x/net/context"
)

type IPlugin interface {
	GetName() string
	GetVersion() string
	GetUsage() string
	GetEventFilter() string
	GetSha1() string
}

type IPluginRegistry interface {
	RegisterPlugin(ctx context.Context, plugin IPlugin) error
}

type Plugin struct {
	*pb.PluginDescriptor
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func NewPluginFromRegisterPluginRequest(pr *pb.RegisterPluginRequest) *Plugin {
	return &Plugin{PluginDescriptor: pr.PluginDescriptor}
}

func (p *Plugin) GetSha1() string {
	val := p.GetName() + p.GetVersion()
	sha1 := sha1.New()
	sha1.Write([]byte(val))
	return fmt.Sprintf("%x", sha1.Sum(nil))
}

type PluginRegistry struct {
}

func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{}
}

func (r *PluginRegistry) RegisterPlugin(ctx context.Context, plugin IPlugin) error {
	return nil
}
