package events

import (
	"fmt"
	"io"
	"runtime"
	"sync"

	pb "github.com/talbright/keds/gen/proto"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

type IEventBus interface {
	AddStream(ctx context.Context, stream pb.KedsService_EventBusServer) (chan struct{}, error)
	Publish(ctx context.Context, event *pb.PluginEvent) error
}

type IListener interface {
	Receive(ctx context.Context, event *pb.PluginEvent) error
	Listen(ctx context.Context) (chan struct{}, error)
}

type StreamListener struct {
	EventBus IEventBus
	Stream   pb.KedsService_EventBusServer
	quitc    chan struct{}
	events   trace.EventLog
}

func NewEventBusMember(eb IEventBus, stream pb.KedsService_EventBusServer) *StreamListener {
	_, file, line, _ := runtime.Caller(1)
	return &StreamListener{
		EventBus: eb,
		Stream:   stream,
		quitc:    make(chan struct{}),
		events:   trace.NewEventLog("plugin.EventBusMember", fmt.Sprintf("%s:%d", file, line)),
	}
}

func (m *StreamListener) Receive(ctx context.Context, event *pb.PluginEvent) (err error) {
	m.events.Printf("forwarding event to plugin: %v", event)
	return m.Stream.Send(event)
}

func (m *StreamListener) Listen(ctx context.Context) (quitc chan struct{}, err error) {
	m.quitc = make(chan struct{})
	go func() {
		for {
			if in, err := m.Stream.Recv(); err == nil {
				m.events.Printf("Listen: event received from plugin: %v", in)
			} else if err == io.EOF {
				m.events.Printf("Listen: EOF")
				close(m.quitc)
				return
			} else {
				m.events.Errorf("Listen: error from stream recv: %v", err)
				close(m.quitc)
				return
			}
		}
	}()
	return m.quitc, err
}

type EventBus struct {
	memberLock *sync.RWMutex
	members    []IListener
	events     trace.EventLog
}

func NewEventBus() *EventBus {
	_, file, line, _ := runtime.Caller(0)
	return &EventBus{
		memberLock: &sync.RWMutex{},
		members:    make([]IListener, 0),
		events:     trace.NewEventLog("plugin.EventBus", fmt.Sprintf("%s:%d", file, line)),
	}
}

func (b *EventBus) AddStream(ctx context.Context, stream pb.KedsService_EventBusServer) (quitc chan struct{}, err error) {
	b.events.Printf("adding new stream")
	member := NewEventBusMember(b, stream)
	if err = b.appendMember(member); err == nil {
		quitc, err = member.Listen(ctx)
	}
	return
}

func (b *EventBus) Publish(ctx context.Context, event *pb.PluginEvent) (err error) {
	b.memberLock.RLock()
	defer b.memberLock.RUnlock()
	b.events.Printf("publishing event '%s' to %d members", event.Name, len(b.members))
	for _, member := range b.members {
		member.Receive(ctx, event)
	}
	return
}

func (b *EventBus) appendMember(member IListener) (err error) {
	b.members = append(b.members, member)
	return
}

func (b *EventBus) deleteMember(member IListener) (err error) {
	return
}
