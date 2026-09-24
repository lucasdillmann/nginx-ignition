package i18n

import (
	"context"
)

func Static(message string) *Message {
	return &Message{
		static: true,
		Key:    message,
	}
}

func M(ctx context.Context, key string) *Message {
	return &Message{
		ctx:       ctx,
		Key:       key,
		Variables: make(map[string]any),
	}
}

func (m Message) V(key string, value any) *Message {
	m.Variables[key] = value
	return &m
}
