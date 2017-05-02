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
)

type IRegistry interface {
	Register(ctx context.Context, plugin IPlugin) error
	Unregister(ctx context.Context) error
	GetFromContext(ctx context.Context) (IPlugin, error)
}

type Registry struct {
	pluginsMutex *sync.RWMutex
	plugins      map[string]IPlugin
}

func NewRegistry() *Registry {
	return &Registry{plugins: make(map[string]IPlugin), pluginsMutex: &sync.RWMutex{}}
}

func (r *Registry) Unregister(ctx context.Context) (err error) {
	var plugin IPlugin
	if plugin, err = r.GetFromContext(ctx); err == nil {
		r.pluginsMutex.Lock()
		defer r.pluginsMutex.Unlock()
		delete(r.plugins, plugin.GetSha1())
	}
	return
}

func (r *Registry) Register(ctx context.Context, plugin IPlugin) (err error) {
	if _, err = r.GetFromContext(ctx); err == ErrPluginMissing {
		err = nil
		r.pluginsMutex.Lock()
		defer r.pluginsMutex.Unlock()
		r.plugins[ut.GetTokenFromContext(ctx)] = plugin
	}
	return
}

func (r *Registry) GetFromContext(ctx context.Context) (plugin IPlugin, err error) {
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

// func (r *Registry) AddCommandForPlugin(plugin IPlugin) (err error) {
// 	if r.rootCmd != nil && plugin.GetRootCommand() != "" {
// 		cmd := &cobra.Command{
// 			Use:                plugin.GetRootCommand(),
// 			Short:              plugin.GetShortDescription(),
// 			Long:               plugin.GetLongDescription(),
// 			DisableFlagParsing: true,
// 			Run: func(cmd *cobra.Command, args []string) {
// 				log.Printf("Cobra.run with args %v", args)
// 				//TODO invoke plugin here
// 			},
// 		}
// 		r.rootCmd.AddCommand(cmd)
// 	}
// 	return
// }
