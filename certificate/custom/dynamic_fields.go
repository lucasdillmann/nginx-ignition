package custom

import "dillmann.com.br/nginx-ignition/core/common/dynamicfields"

const (
	textFieldUploadModeID = "textField"
	fileUploadModeID      = "fileUpload"
)

var (
	uploadModeField = dynamicfields.DynamicField{
		ID:          "uploadMode",
		Priority:    0,
		Description: "Upload mode",
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.EnumType,
		EnumOptions: []dynamicfields.EnumOption{
			{
				ID:          textFieldUploadModeID,
				Description: "PEM-encoded text",
			},
			{
				ID:          fileUploadModeID,
				Description: "PEM-encoded file",
			},
		},
	}

	publicKeyTextField = dynamicfields.DynamicField{
		ID:          "publicKeyPem",
		Priority:    1,
		Description: "Public key",
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.MultiLineTextType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       textFieldUploadModeID,
		}},
	}

	privateKeyTextField = dynamicfields.DynamicField{
		ID:          "privateKeyPem",
		Priority:    2,
		Description: "Private key",
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.MultiLineTextType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       textFieldUploadModeID,
		}},
	}

	certificationChainTextField = dynamicfields.DynamicField{
		ID:          "certificationChainPem",
		Priority:    3,
		Description: "Certification chain",
		Required:    false,
		Sensitive:   true,
		Type:        dynamicfields.MultiLineTextType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       textFieldUploadModeID,
		}},
	}

	publicKeyFileField = dynamicfields.DynamicField{
		ID:          "publicKeyFile",
		Priority:    1,
		Description: "Public key",
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.FileType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       fileUploadModeID,
		}},
	}

	privateKeyFileField = dynamicfields.DynamicField{
		ID:          "privateKeyFile",
		Priority:    2,
		Description: "Private key",
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.FileType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       fileUploadModeID,
		}},
	}

	certificationChainFileField = dynamicfields.DynamicField{
		ID:          "certificationChainFile",
		Priority:    3,
		Description: "Certification chain",
		Required:    false,
		Sensitive:   true,
		Type:        dynamicfields.FileType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       fileUploadModeID,
		}},
	}
)
