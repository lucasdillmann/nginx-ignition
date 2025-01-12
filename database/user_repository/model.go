package user_repository

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type userModel struct {
	bun.BaseModel `bun:"user"`

	ID           uuid.UUID `bun:"id,pk"`
	Enabled      bool      `bun:"enabled,notnull"`
	Name         string    `bun:"name,notnull"`
	Username     string    `bun:"username,notnull"`
	PasswordHash string    `bun:"password_hash,notnull"`
	PasswordSalt string    `bun:"password_salt,notnull"`
	Role         string    `bun:"role,notnull"`
}
