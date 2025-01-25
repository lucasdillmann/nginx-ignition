package integration

import (
	"github.com/uptrace/bun"
)

type integrationModel struct {
	bun.BaseModel `bun:"integration"`

	ID         string `bun:"id,pk"`
	Enabled    bool   `bun:"enabled,notnull"`
	Parameters string `bun:"parameters,notnull"`
}
