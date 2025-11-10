package custom

import "dillmann.com.br/nginx-ignition/core/common/dynamicfields"

var (
	publicKeyField = dynamicfields.DynamicField{
		ID:          "publicKey",
		Priority:    0,
		Description: "Certificate file (PEM encoded) with the public key",
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.FileType,
	}

	privateKeyField = dynamicfields.DynamicField{
		ID:          "privateKey",
		Priority:    1,
		Description: "Certificate file (PEM encoded) with the private key",
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.FileType,
	}

	certificationChainField = dynamicfields.DynamicField{
		ID:          "certificationChain",
		Priority:    2,
		Description: "Certification chain file (PEM encoded)",
		Required:    false,
		Sensitive:   true,
		Type:        dynamicfields.FileType,
	}
)
