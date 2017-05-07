package events_test

import (
	. "github.com/talbright/keds/events"
	pb "github.com/talbright/keds/gen/proto"
	"github.com/talbright/keds/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Event", func() {
	var plug *plugin.Plugin
	var descriptor *pb.PluginDescriptor
	BeforeEach(func() {
		descriptor = &pb.PluginDescriptor{
			Name:             "testPlugin",
			Usage:            "testPlugin",
			EventFilter:      "*",
			Version:          "1",
			RootCommand:      "testPlugin",
			ShortDescription: "This is a testPlugin",
			LongDescription:  "This is a testPlugin",
		}
		plug = plugin.NewPlugin(descriptor)
	})
	Context("NewEvent", func() {
		It("should set defaults", func() {
			e := NewEvent(nil, nil)
			Expect(e.Args).ShouldNot(BeNil())
			Expect(e.Data).ShouldNot(BeNil())
			Expect(e.PluginEvent).ShouldNot(BeNil())
			e = NewEvent(nil, plug)
			Expect(e.Source).Should(Equal(plug.GetName()))
		})
	})
	Context("NewEventWithOptions", func() {
		It("should apply options to constructor", func() {
			args := []string{"args"}
			data := map[string]string{"data": "data"}
			e := NewEventWithOptions(
				WithName("foo"),
				WithSourcePlugin(plug),
				WithTarget("bar"),
				WithArgs(args),
				WithData(data),
				WithExitCode(1))
			Expect(e.GetName()).Should(Equal("foo"))
			Expect(e.GetTarget()).Should(Equal("bar"))
			Expect(e.SourcePlugin).Should(Equal(plug))
			Expect(e.GetArgs()).Should(ContainElement("args"))
			Expect(e.GetData()).Should(ContainElement("data"))
			Expect(e.GetData()["exit_code"]).Should(Equal("1"))
		})
	})
	Context("CreateEventServerQuit", func() {
		It("should create the event", func() {
			e := CreateEventServerQuit(plug, 2)
			Expect(e.GetName()).Should(Equal("keds.exit"))
			Expect(e.GetData()["exit_code"]).Should(Equal("2"))
			Expect(e.GetSource()).Should(Equal("testPlugin"))
		})
	})
	Context("CreateEventCommandInvoked", func() {
		It("should create the event", func() {
			args := []string{"abc"}
			e := CreateEventCommandInvoked(plug, args)
			Expect(e.GetName()).Should(Equal("keds.command_invoked"))
			Expect(e.GetArgs()).Should(ContainElement("abc"))
			Expect(e.GetTarget()).Should(Equal("testPlugin"))
		})
	})
})
