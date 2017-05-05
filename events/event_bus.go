package events

import (
	"fmt"
	"runtime"
	"sync"

	pb "github.com/talbright/keds/gen/proto"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

type IEventBus interface {
	AddListener(ctx context.Context, listener IListener) (chan struct{}, error)
	Publish(ctx context.Context, event *pb.PluginEvent) error
}

type IListener interface {
	Receive(ctx context.Context, event *pb.PluginEvent) error
	Listen(ctx context.Context) (chan struct{}, error)
}

type EventBus struct {
	memberLock *sync.RWMutex
	listeners  []IListener
	events     trace.EventLog
}

func NewEventBus() *EventBus {
	_, file, line, _ := runtime.Caller(0)
	return &EventBus{
		memberLock: &sync.RWMutex{},
		listeners:  make([]IListener, 0),
		events:     trace.NewEventLog("plugin.EventBus", fmt.Sprintf("%s:%d", file, line)),
	}
}

func (b *EventBus) AddListener(ctx context.Context, listener IListener) (quitc chan struct{}, err error) {
	b.memberLock.RLock()
	defer b.memberLock.RUnlock()
	b.listeners = append(b.listeners, listener)
	quitc, err = listener.Listen(ctx)
	return
}

func (b *EventBus) Publish(ctx context.Context, event *pb.PluginEvent) (err error) {
	b.memberLock.RLock()
	defer b.memberLock.RUnlock()
	b.events.Printf("publishing event '%s' to %d listeners", event.Name, len(b.listeners))
	for _, member := range b.listeners {
		member.Receive(ctx, event)
	}
	return
}
