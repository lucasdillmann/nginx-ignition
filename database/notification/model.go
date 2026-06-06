package notification

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type notificationModel struct {
	bun.BaseModel     `bun:"notification"`
	CreatedAt         time.Time  `bun:"created_at,notnull"`
	ReadAt            *time.Time `bun:"read_at"`
	Payload           string     `bun:"payload,notnull"`
	Title             string     `bun:"title,notnull"`
	Summary           string     `bun:"summary,notnull"`
	Category          string     `bun:"category,notnull"`
	ID                uuid.UUID  `bun:"id,pk"`
	UserID            uuid.UUID  `bun:"user_id,notnull"`
	DeliveryCompleted bool       `bun:"delivery_completed,notnull"`
}

type configurationModel struct {
	bun.BaseModel `bun:"notification_configuration"`
	Categories    *string   `bun:"categories"`
	Parameters    string    `bun:"parameters,notnull"`
	Name          string    `bun:"name,notnull"`
	Provider      string    `bun:"provider,notnull"`
	ID            uuid.UUID `bun:"id,pk"`
	UserID        uuid.UUID `bun:"user_id,notnull"`
	Enabled       bool      `bun:"enabled,notnull"`
}

type providerSubmissionModel struct {
	bun.BaseModel   `bun:"notification_provider_submission"`
	LastError       *string    `bun:"last_error"`
	LastAttemptAt   *time.Time `bun:"last_attempt_at"`
	SucceededAt     *time.Time `bun:"succeeded_at"`
	Status          string     `bun:"status,notnull"`
	Provider        string     `bun:"provider,notnull"`
	AttemptCount    int        `bun:"attempt_count,notnull"`
	ID              uuid.UUID  `bun:"id,pk"`
	NotificationID  uuid.UUID  `bun:"notification_id,notnull"`
	ConfigurationID uuid.UUID  `bun:"configuration_id,notnull"`
}

type relatedEntityModel struct {
	bun.BaseModel  `bun:"notification_related_entity"`
	Name           *string   `bun:"name"`
	EntityType     string    `bun:"entity_type,notnull"`
	NotificationID uuid.UUID `bun:"notification_id,pk"`
	EntityID       uuid.UUID `bun:"entity_id,pk"`
}
