package plugin

import (
	ut "github.com/talbright/keds/utils/token"
	"golang.org/x/net/context"

	"fmt"
	"sync"
)

var (
	ErrPluginExists       = fmt.Errorf("plugin already registered")
	ErrPluginTokenMissing = fmt.Errorf("plugin missing token")
	ErrPluginMissing      = fmt.Errorf("plugin missing")
	registry              = NewRegistry()
)

func DefaultRegistry() IRegistry { return registry }

type IRegistry interface {
	Register(ctx context.Context, plugin *Plugin) error
	Unregister(ctx context.Context) error
	GetFromContext(ctx context.Context) (*Plugin, error)
}

type Registry struct {
	pluginsMutex *sync.RWMutex
	plugins      map[string]*Plugin
}

func NewRegistry() *Registry {
	return &Registry{plugins: make(map[string]*Plugin), pluginsMutex: &sync.RWMutex{}}
}

func (r *Registry) Unregister(ctx context.Context) (err error) {
	var plugin *Plugin
	if plugin, err = r.GetFromContext(ctx); err == nil {
		r.pluginsMutex.Lock()
		defer r.pluginsMutex.Unlock()
		delete(r.plugins, plugin.GetSha1())
	}
	return
}

func (r *Registry) Register(ctx context.Context, plugin *Plugin) (err error) {
	if _, err = r.GetFromContext(ctx); err == ErrPluginMissing {
		err = nil
		r.pluginsMutex.Lock()
		defer r.pluginsMutex.Unlock()
		r.plugins[ut.GetTokenFromContext(ctx)] = plugin
	}
	return
}

func (r *Registry) GetFromContext(ctx context.Context) (plugin *Plugin, err error) {
	r.pluginsMutex.RLock()
	defer r.pluginsMutex.RUnlock()
	token := ut.GetTokenFromContext(ctx)
	if token == "" {
		return nil, ErrPluginTokenMissing
	}
	var ok bool
	if plugin, ok = r.plugins[token]; !ok {
		err = ErrPluginMissing
	}
	return
}
