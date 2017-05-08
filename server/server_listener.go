package server

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

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
	if event.GetSource() != "keds" {
		m.eventc <- event
	}
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
	if event.Name == "keds.exit" {
		exitCode := 0
		if strCode, ok := event.Data["exit_code"]; ok {
			exitCode, _ = strconv.Atoi(strCode)
		}
		os.Exit(exitCode)
	}
}

func (m ServerListener) String() string {
	return fmt.Sprintf("server.ServerListener (internal)")
}
