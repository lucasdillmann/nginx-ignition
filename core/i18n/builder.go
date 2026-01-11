package i18n

import (
	"context"
)

func K(ctx context.Context, key string) *Message {
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
