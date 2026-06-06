package notification

import (
	"context"
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type ContentSection struct {
	Title *i18n.DetachedMessage
	Body  i18n.DetachedMessage
}

type Action struct {
	Label i18n.DetachedMessage
	URL   string
}

type RelatedEntity struct {
	Type string
	Name string
	ID   uuid.UUID
}

type SendRequest struct {
	Title           i18n.DetachedMessage
	Summary         i18n.DetachedMessage
	OccurredAt      time.Time
	Tags            map[string]string
	Category        Category
	Sections        []ContentSection
	Actions         []Action
	RelatedEntities []RelatedEntity
}

type Deliverable struct {
	OccurredAt time.Time
	Tags       map[string]string
	Title      string
	Summary    string
	Category   Category
	Sections   []DeliverableContentSection
	Actions    []DeliverableAction
}

type Provider interface {
	ID() string
	Name(ctx context.Context) *i18n.Message
	ImportantInstructions(ctx context.Context) []*i18n.Message
	ConfigurationFields(ctx context.Context) []dynamicfields.DynamicField
	Send(ctx context.Context, parameters map[string]any, deliverable Deliverable) error
}
