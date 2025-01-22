package custom

import "dillmann.com.br/nginx-ignition/core/common/dynamic_fields"

var (
	publicKey = dynamic_fields.DynamicField{
		ID:          "publicKey",
		Priority:    0,
		Description: "Certificate file (PEM encoded) with the public key",
		Required:    true,
		Sensitive:   true,
		Type:        dynamic_fields.FileType,
	}

	privateKey = dynamic_fields.DynamicField{
		ID:          "privateKey",
		Priority:    1,
		Description: "Certificate file (PEM encoded) with the private key",
		Required:    true,
		Sensitive:   true,
		Type:        dynamic_fields.FileType,
	}

	certificationChain = dynamic_fields.DynamicField{
		ID:          "certificationChain",
		Priority:    2,
		Description: "Certification chain file (PEM encoded)",
		Required:    false,
		Sensitive:   true,
		Type:        dynamic_fields.FileType,
	}
)
