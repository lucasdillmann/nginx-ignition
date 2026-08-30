package integration

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type integrationModel struct {
	bun.BaseModel `bun:"integration"`

	Driver     string    `bun:"driver,notnull"`
	Name       string    `bun:"name,notnull"`
	Parameters string    `bun:"parameters,notnull"`
	ID         uuid.UUID `bun:"id,pk"`
	Enabled    bool      `bun:"enabled,notnull"`
}
