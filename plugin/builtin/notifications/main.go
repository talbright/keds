package main

import (
	"log"

	"sync"

	"github.com/spf13/cobra"
	"github.com/talbright/keds/client"
	pb "github.com/talbright/keds/gen/proto"
	"golang.org/x/net/context"
)

type Notifier interface {
	Notify(ctx context.Context, cli *client.Client, event *pb.PluginEvent)
}

type NotificationHandler struct {
	mutex     *sync.RWMutex
	notifiers []Notifier
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{notifiers: make([]Notifier, 0), mutex: &sync.RWMutex{}}
}

func (e *NotificationHandler) AddNotifier(notifier Notifier) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.notifiers = append(e.notifiers, notifier)
}

func (e *NotificationHandler) OnCommandInvoked(cli *client.Client, event *pb.PluginEvent, cmd *cobra.Command, args []string) {
	log.Println("NotificationHandler.OnCommandInvoked")
}

func (e *NotificationHandler) OnBusEvent(cli *client.Client, event *pb.PluginEvent) {
	log.Println("NotificationHandler.OnBusEvent")
}

func (e *NotificationHandler) OnInitCommand(cli *client.Client, cmd *cobra.Command) {
	log.Println("NotificationHandler.OnInitCommand")
}

func (e *NotificationHandler) OnQuit(cli *client.Client) {
	log.Println("NotificationHandler.OnQuit")
}

//TODO
//- make sure acces to client properties are threadsafe
//- create deep copy of event
//- consider that notifications may be out of received by the notifier out of order
func (e *NotificationHandler) dispatch(cli *client.Client, event *pb.PluginEvent) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	for _, n := range e.notifiers {
		go func() {
			ctx := context.Background()
			n.Notify(ctx, cli, event)
		}()
	}
}

func main() {
	handler := &NotificationHandler{}
	notifications, err := client.NewClient(handler)
	if err != nil {
		log.Fatalf("client error: %v", err)
	}
	notifications.Run()
}
