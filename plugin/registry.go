package plugin

import (
	"log"

	"github.com/spf13/cobra"
	ut "github.com/talbright/keds/utils/token"
	"golang.org/x/net/context"

	"fmt"
	"sync"
)

var (
	registry              = NewPluginRegistry(nil)
	ErrPluginExists       = fmt.Errorf("plugin already registered")
	ErrPluginTokenMissing = fmt.Errorf("plugin missing token")
	ErrPluginMissing      = fmt.Errorf("plugin missing")
)

//TODO remove
func DefaultRegistry() IPluginRegistry { return registry }

type IPluginRegistry interface {
	RegisterPlugin(ctx context.Context, plugin IPlugin) error
	UnRegisterPlugin(ctx context.Context) error
	GetPluginFromContext(ctx context.Context) (IPlugin, error)
}

type PluginRegistry struct {
	pluginsMutex *sync.RWMutex
	plugins      map[string]IPlugin
	rootCmd      *cobra.Command
}

func NewPluginRegistry(rootCmd *cobra.Command) *PluginRegistry {
	return &PluginRegistry{plugins: make(map[string]IPlugin), pluginsMutex: &sync.RWMutex{}, rootCmd: rootCmd}
}

func (r *PluginRegistry) UnRegisterPlugin(ctx context.Context) (err error) {
	var plugin IPlugin
	if plugin, err = r.GetPluginFromContext(ctx); err == nil {
		r.pluginsMutex.Lock()
		defer r.pluginsMutex.Unlock()
		delete(r.plugins, plugin.GetSha1())
	}
	return
}

func (r *PluginRegistry) RegisterPlugin(ctx context.Context, plugin IPlugin) (err error) {
	if _, err = r.GetPluginFromContext(ctx); err == ErrPluginMissing {
		err = nil
		r.pluginsMutex.Lock()
		defer r.pluginsMutex.Unlock()
		r.plugins[ut.GetTokenFromContext(ctx)] = plugin
		r.AddCommandForPlugin(plugin)
	}
	return
}

func (r *PluginRegistry) GetPluginFromContext(ctx context.Context) (IPlugin, error) {
	r.pluginsMutex.RLock()
	defer r.pluginsMutex.RUnlock()
	token := ut.GetTokenFromContext(ctx)
	if token == "" {
		return nil, ErrPluginTokenMissing
	}
	if plugin, ok := r.plugins[token]; !ok {
		return nil, ErrPluginMissing
	} else {
		return plugin, nil
	}
}

func (r *PluginRegistry) AddCommandForPlugin(plugin IPlugin) (err error) {
	if r.rootCmd != nil {
		cmd := &cobra.Command{
			Use:   plugin.GetName(),
			Short: plugin.GetUsage(),
			Long:  plugin.GetUsage(),
			Run: func(cmd *cobra.Command, args []string) {
				log.Printf("Cobra.run with args %v", args)
			},
		}
		r.rootCmd.AddCommand(cmd)
	}
	return
}
