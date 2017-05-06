package server

import (
	"fmt"
	"io"
	"runtime"

	"github.com/talbright/keds/events"
	pb "github.com/talbright/keds/gen/proto"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

type StreamListener struct {
	EventBus events.IEventBus
	Stream   pb.KedsService_EventBusServer
	quitc    chan struct{}
	events   trace.EventLog
}

func NewStreamListener(eb events.IEventBus, stream pb.KedsService_EventBusServer) *StreamListener {
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
	if event.Source != "example" {
		err = m.Stream.Send(event)
	}
	return
}

func (m *StreamListener) Listen(ctx context.Context) (quitc chan struct{}, err error) {
	m.quitc = make(chan struct{})
	go func() {
		for {
			if in, err := m.Stream.Recv(); err == nil {
				m.events.Printf("Listen: event received from plugin: %v", in)
				m.EventBus.Publish(ctx, in)
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
