package server

import (
	"fmt"
	"os"
	"runtime"

	"github.com/talbright/keds/events"
	pb "github.com/talbright/keds/gen/proto"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

type ServerListener struct {
	EventBus events.IEventBus
	quitc    chan struct{}
	events   trace.EventLog
	server   *KedsRPCServer
	eventc   chan *pb.PluginEvent
}

func NewServerListener(eb events.IEventBus, server *KedsRPCServer) *ServerListener {
	_, file, line, _ := runtime.Caller(1)
	return &ServerListener{
		EventBus: eb,
		server:   server,
		quitc:    make(chan struct{}),
		eventc:   make(chan *pb.PluginEvent),
		events:   trace.NewEventLog("server.ServerListener", fmt.Sprintf("%s:%d", file, line)),
	}
}

func (m *ServerListener) Receive(ctx context.Context, event *pb.PluginEvent) (err error) {
	m.events.Printf("received event: %v", event)
	m.eventc <- event
	return
}

func (m *ServerListener) Listen(ctx context.Context) (quitc chan struct{}, err error) {
	m.quitc = make(chan struct{})
	go func() {
		for event := range m.eventc {
			m.handleEvent(event)
		}
	}()
	return m.quitc, err
}

func (m *ServerListener) handleEvent(event *pb.PluginEvent) {
	if event.Name == "server.quit" {
		os.Exit(0)
	}
}
