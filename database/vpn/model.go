package vpn

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type vpnModel struct {
	bun.BaseModel `bun:"vpn"`

	ID         uuid.UUID `bun:"id,pk"`
	Driver     string    `bun:"driver,notnull"`
	Name       string    `bun:"name,notnull"`
	Enabled    bool      `bun:"enabled,notnull"`
	Parameters string    `bun:"parameters,notnull"`
}
