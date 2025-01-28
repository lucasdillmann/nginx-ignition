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
