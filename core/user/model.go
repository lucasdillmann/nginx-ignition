package user

import "github.com/google/uuid"

type SaveRequest struct {
	ID       uuid.UUID
	Enabled  bool
	Name     string
	Username string
	Password *string
	Role     Role
}

type Role string

const (
	RegularRole = Role("REGULAR_USER")
	AdminRole   = Role("ADMIN")
)

type User struct {
	ID           uuid.UUID
	Enabled      bool
	Name         string
	Username     string
	PasswordHash string
	PasswordSalt string
	Role         Role
}
