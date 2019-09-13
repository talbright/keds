package main

import (
	"mime"
	"strings"

	"github.com/talbright/keds/gen/proto"
)

func init() {
	mime.AddExtensionType(".txt", "text/plain")
}

//Content models the data and mime type (from an event.)
type Content struct {
	Data string
	Type string
}

//Message contains all the content extracted from an event.
type Message struct {
	content []Content
}

//NewMessage creates a new notification message extracted from an event.
func NewMessage(event *proto.PluginEvent) *Message {
	m := &Message{content: make([]Content, 0)}
	m.extractContent(event)
	return m
}

//GetData finds the notification data by mime type.
func (m *Message) GetData(mimeType string) (data string) {
	for _, v := range m.content {
		if v.Type == mimeType {
			return v.Data
		}
	}
	return ""
}

func (m *Message) extractContent(event *proto.PluginEvent) {
	for k, v := range event.Data {
		parts := strings.Split(k, ".")
		if len(parts) == 2 && parts[0] == "message" {
			if mimeType := mime.TypeByExtension("." + parts[1]); mimeType != "" {
				m.content = append(m.content, Content{Data: v, Type: mimeType})
			}
		}
	}
}
