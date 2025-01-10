package host

import "github.com/google/uuid"

type BindingType string

const (
	HttpBindingType  = BindingType("HTTP")
	HttpsBindingType = BindingType("HTTPS")
)

type Binding struct {
	ID            uuid.UUID
	Type          BindingType
	IP            string
	Port          int
	CertificateID *uuid.UUID
}
