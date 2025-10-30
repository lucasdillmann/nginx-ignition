package integration

import (
	"github.com/google/uuid"
)

type Integration struct {
	ID         uuid.UUID
	Driver     string
	Name       string
	Enabled    bool
	Parameters map[string]any
}
