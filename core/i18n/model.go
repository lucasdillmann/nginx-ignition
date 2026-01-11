package i18n

import (
	"context"
	"encoding/json"

	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

type Message struct {
	ctx       context.Context
	Variables map[string]any
	Key       string
}

type Dictionary struct {
	Language  language.Tag
	Templates map[string]string
}

func (m Message) String() string {
	commands := container.Get[Commands]()

	var lang *language.Tag
	if ctxLang, casted := m.ctx.Value(ContextKey).(language.Tag); casted {
		lang = &ctxLang
	} else {
		lang = ptr.Of(commands.DefaultLanguage())
		log.Warnf("Language not found in context. Using %s as fallback.", lang.String())
	}

	return commands.Translate(*lang, m.Key, m.Variables)
}

func (m Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}
