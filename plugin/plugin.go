package plugin

import (
	pb "github.com/talbright/keds/gen/proto"

	"crypto/sha1"
	"fmt"
)

type IPlugin interface {
	GetName() string
	GetVersion() string
	GetUsage() string
	GetEventFilter() string
	GetSha1() string
	GetSha1Short() string
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

func (p *Plugin) GetSha1Short() string {
	return fmt.Sprintf("%.7s", p.GetSha1())
}

func (p Plugin) String() string {
	return fmt.Sprintf("plugin '%s' v%s (%s)", p.GetName(), p.GetVersion(), p.GetSha1Short())
}
