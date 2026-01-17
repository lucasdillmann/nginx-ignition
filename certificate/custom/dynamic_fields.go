package custom

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	textFieldUploadModeID = "textField"
	fileUploadModeID      = "fileUpload"
	uploadModeFieldID     = "uploadMode"
	publicKeyTextFieldID  = "publicKeyPem"
	privateKeyTextFieldID = "privateKeyPem"
	chainTextFieldID      = "certificationChainPem"
	publicKeyFileFieldID  = "publicKeyFile"
	privateKeyFileFieldID = "privateKeyFile"
	chainFileFieldID      = "certificationChainFile"
)

func dynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	uploadModeField := dynamicfields.DynamicField{
		ID:          uploadModeFieldID,
		Priority:    0,
		Description: i18n.M(ctx, i18n.K.CertificateCommonCustomUploadMode),
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.EnumType,
		EnumOptions: []dynamicfields.EnumOption{
			{
				ID:          textFieldUploadModeID,
				Description: i18n.M(ctx, i18n.K.CertificateCommonCustomUploadModeText),
			},
			{
				ID:          fileUploadModeID,
				Description: i18n.M(ctx, i18n.K.CertificateCommonCustomUploadModeFile),
			},
		},
	}

	publicKeyTextField := dynamicfields.DynamicField{
		ID:          publicKeyTextFieldID,
		Priority:    1,
		Description: i18n.M(ctx, i18n.K.CertificateCommonCustomPublicKey),
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.MultiLineTextType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       textFieldUploadModeID,
		}},
	}

	privateKeyTextField := dynamicfields.DynamicField{
		ID:          privateKeyTextFieldID,
		Priority:    2,
		Description: i18n.M(ctx, i18n.K.CertificateCommonCustomPrivateKey),
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.MultiLineTextType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       textFieldUploadModeID,
		}},
	}

	certificationChainTextField := dynamicfields.DynamicField{
		ID:          chainTextFieldID,
		Priority:    3,
		Description: i18n.M(ctx, i18n.K.CertificateCommonCustomChain),
		Required:    false,
		Sensitive:   true,
		Type:        dynamicfields.MultiLineTextType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       textFieldUploadModeID,
		}},
	}

	publicKeyFileField := dynamicfields.DynamicField{
		ID:          publicKeyFileFieldID,
		Priority:    1,
		Description: i18n.M(ctx, i18n.K.CertificateCommonCustomPublicKey),
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.FileType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       fileUploadModeID,
		}},
	}

	privateKeyFileField := dynamicfields.DynamicField{
		ID:          privateKeyFileFieldID,
		Priority:    2,
		Description: i18n.M(ctx, i18n.K.CertificateCommonCustomPrivateKey),
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.FileType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       fileUploadModeID,
		}},
	}

	certificationChainFileField := dynamicfields.DynamicField{
		ID:          chainFileFieldID,
		Priority:    3,
		Description: i18n.M(ctx, i18n.K.CertificateCommonCustomChain),
		Required:    false,
		Sensitive:   true,
		Type:        dynamicfields.FileType,
		Conditions: []dynamicfields.Condition{{
			ParentField: uploadModeField.ID,
			Value:       fileUploadModeID,
		}},
	}

	return []dynamicfields.DynamicField{
		uploadModeField,
		publicKeyTextField,
		privateKeyTextField,
		certificationChainTextField,
		publicKeyFileField,
		privateKeyFileField,
		certificationChainFileField,
	}
}
