package binding

import (
	"github.com/google/uuid"
)

type Type string

const (
	HTTPBindingType  Type = "HTTP"
	HTTPSBindingType Type = "HTTPS"
)

type Binding struct {
	CertificateID *uuid.UUID
	Type          Type
	IP            string
	Port          int
	ID            uuid.UUID
}
