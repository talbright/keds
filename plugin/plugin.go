package plugin

import (
	pb "github.com/talbright/keds/gen/proto"

	"crypto/sha1"
	"fmt"
)

type Plugin struct {
	*pb.PluginDescriptor
}

func NewPlugin(descriptor *pb.PluginDescriptor) *Plugin {
	return &Plugin{PluginDescriptor: descriptor}
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
	return fmt.Sprintf("'%s' v%s (%s)", p.GetName(), p.GetVersion(), p.GetSha1Short())
}

func (p *Plugin) Proto() *pb.PluginDescriptor {
	return p.PluginDescriptor
}
