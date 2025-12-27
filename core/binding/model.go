package binding

import (
	"github.com/google/uuid"
)

type Type string

const (
	HttpBindingType  Type = "HTTP"
	HttpsBindingType Type = "HTTPS"
)

type Binding struct {
	CertificateID *uuid.UUID
	Type          Type
	IP            string
	Port          int
	ID            uuid.UUID
}
