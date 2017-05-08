package server

import (
	"fmt"
	"io"
	"runtime"

	"github.com/talbright/keds/events"
	pb "github.com/talbright/keds/gen/proto"
	"github.com/talbright/keds/plugin"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

type StreamListener struct {
	EventBus events.IEventBus
	Stream   pb.KedsService_EventBusServer
	quitc    chan struct{}
	events   trace.EventLog
	listener *plugin.Plugin
	name     string
}

func NewStreamListener(eb events.IEventBus, stream pb.KedsService_EventBusServer, name string) *StreamListener {
	_, file, line, _ := runtime.Caller(1)
	return &StreamListener{
		EventBus: eb,
		Stream:   stream,
		quitc:    make(chan struct{}),
		name:     name,
		events:   trace.NewEventLog("plugin.StreamListener", fmt.Sprintf("%s:%d", file, line)),
	}
}

func (m *StreamListener) Receive(ctx context.Context, event *pb.PluginEvent) (err error) {
	var sender *plugin.Plugin
	if sender, err = plugin.DefaultRegistry().GetFromContext(ctx); err == nil {
		if sender.GetSha1() != m.listener.GetSha1() {
			err = m.Stream.Send(event)
		}
	} else if err == plugin.ErrPluginTokenMissing && event.GetSource() == "keds" {
		err = m.Stream.Send(event)
	}
	return
}

func (m *StreamListener) Listen(ctx context.Context) (quitc chan struct{}, err error) {
	m.quitc = make(chan struct{})
	if m.listener, err = plugin.DefaultRegistry().GetFromContext(ctx); err != nil {
		return
	}
	go func() {
		for {
			if in, err := m.Stream.Recv(); err == nil {
				m.events.Printf("Listen: event received from plugin: %v", in)
				in.Source = in.GetName()
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

func (m StreamListener) String() string {
	return fmt.Sprintf("server.StreamListener (%s)", m.name)
}
