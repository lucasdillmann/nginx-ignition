package smtp

import (
	"context"
)

type stubMailSender struct {
	lastMessage    []byte
	lastRecipients []string
	lastSettings   mailSettings
}

func (stub *stubMailSender) Send(
	_ context.Context,
	settings mailSettings,
	recipients []string,
	message []byte,
) error {
	stub.lastSettings = settings
	stub.lastRecipients = append([]string(nil), recipients...)
	stub.lastMessage = append([]byte(nil), message...)
	return nil
}
