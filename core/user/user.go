package user

import (
	"github.com/google/uuid"
)

type Role string

const (
	RegularRole = Role("REGULAR")
	AdminRole   = Role("ADMIN")
)

type User struct {
	Id           uuid.UUID
	Enabled      bool
	Name         string
	Username     string
	PasswordHash string
	PasswordSalt string
	Role         Role
}
