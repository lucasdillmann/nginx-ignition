package integration

import (
	"github.com/google/uuid"
)

type Integration struct {
	Parameters map[string]any
	Driver     string
	Name       string
	ID         uuid.UUID
	Enabled    bool
}
