package notification

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamicfield"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type notificationResponse struct {
	CreatedAt       time.Time                `json:"createdAt"`
	OccurredAt      time.Time                `json:"occurredAt"`
	Tags            map[string]string        `json:"tags"`
	Title           string                   `json:"title"`
	Summary         string                   `json:"summary"`
	Category        string                   `json:"category"`
	Sections        []contentSectionResponse `json:"sections"`
	RelatedEntities []relatedEntityResponse  `json:"relatedEntities"`
	Actions         []actionResponse         `json:"actions"`
	ID              uuid.UUID                `json:"id"`
	Read            bool                     `json:"read"`
}

type contentSectionResponse struct {
	Title *string `json:"title,omitempty"`
	Body  string  `json:"body"`
}

type actionResponse struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type relatedEntityResponse struct {
	Name string    `json:"name,omitempty"`
	Type string    `json:"type"`
	ID   uuid.UUID `json:"id"`
}

type unreadCountResponse struct {
	Count int `json:"count"`
}

type categoryResponse struct {
	Name *i18n.Message `json:"name"`
	ID   string        `json:"id"`
}

type configurationRequest struct {
	Parameters map[string]any  `json:"parameters"`
	Name       string          `json:"name"`
	Provider   string          `json:"provider"`
	Categories json.RawMessage `json:"categories"`
	Enabled    bool            `json:"enabled"`
}

type configurationResponse struct {
	Parameters map[string]any  `json:"parameters"`
	Name       string          `json:"name"`
	Provider   string          `json:"provider"`
	Categories json.RawMessage `json:"categories"`
	ID         uuid.UUID       `json:"id"`
	Enabled    bool            `json:"enabled"`
}

type availableProviderResponse struct {
	Name                  *i18n.Message           `json:"name"`
	ID                    string                  `json:"id"`
	ImportantInstructions []*i18n.Message         `json:"importantInstructions"`
	ConfigurationFields   []dynamicfield.Response `json:"configurationFields"`
}
