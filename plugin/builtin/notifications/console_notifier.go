package main

import (
	"mime"
	"strings"

	"github.com/talbright/keds/gen/proto"
)

func init() {
	mime.AddExtensionType(".txt", "text/plain")
}

type Content struct {
	Data string
	Type string
}

type Message struct {
	content []Content
}

func NewMessage(event *proto.PluginEvent) *Message {
	m := &Message{content: make([]Content, 0)}
	m.extractContent(event)
	return m
}

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
