package i18n

import (
	"context"
	"encoding/json"

	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

type Message struct {
	ctx       context.Context
	Variables map[string]any
	Key       string
	static    bool
}

func (m Message) String() string {
	if m.static || m.ctx == nil || !container.Ready() {
		return m.Key
	}

	commands := container.Get[Commands]()
	if commands == nil {
		return m.Key
	}

	var lang *language.Tag
	if ctxLang, casted := m.ctx.Value(ContextKey).(language.Tag); casted {
		lang = &ctxLang
	} else {
		lang = ptr.Of(commands.DefaultLanguage())
	}

	return commands.Translate(*lang, m.Key, m.Variables)
}

func (m Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}
