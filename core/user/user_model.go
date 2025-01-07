package user

import "github.com/google/uuid"

type SaveRequest struct {
	Id       uuid.UUID
	Enabled  bool
	Name     string
	Username string
	Password string
	Role     Role
}
