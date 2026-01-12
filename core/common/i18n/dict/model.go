package dict

import (
	"golang.org/x/text/language"
)

type Dictionary interface {
	Language() language.Tag
	Translate(messageKey string, variables map[string]any) string
	Templates() map[string]string
}

type Messages struct {
	AccessListValidationInvalidAddress             string
	BindingValidationCertificateIDNotAllowed       string
	BindingValidationCertificateIDNotFound         string
	BindingValidationCertificateIDRequired         string
	BindingValidationInvalidIP                     string
	BindingValidationInvalidType                   string
	CacheValidationExtensionDotNotAllowed          string
	CacheValidationInvalidMethod                   string
	CacheValidationInvalidStaleOption              string
	CacheValidationInvalidStatusCode               string
	CommonValidationAbsolutePathRequired           string
	CommonValidationAtLeastOneRequired             string
	CommonValidationBetweenValues                  string
	CommonValidationCannotBeEmpty                  string
	CommonValidationCannotBeNegative               string
	CommonValidationCannotBeZero                   string
	CommonValidationDuplicatedValue                string
	CommonValidationInvalidDomainName              string
	CommonValidationInvalidURI                     string
	CommonValidationInvalidURL                     string
	CommonValidationInvalidValue                   string
	CommonValidationStartsWithSlashRequired        string
	CommonValidationTooLong                        string
	CommonValidationTooShort                       string
	CommonValidationValueMissing                   string
	DynamicFieldValidationInvalidBoolean           string
	DynamicFieldValidationInvalidEmail             string
	DynamicFieldValidationInvalidFileEncodedBase64 string
	DynamicFieldValidationInvalidText              string
	DynamicFieldValidationNotRecognizedOption      string
	DynamicFieldValidationUnknownFieldType         string
	HostValidationAccessListNotFound               string
	HostValidationBindingsMustBeEmptyForGlobal     string
	HostValidationCacheNotFound                    string
	HostValidationDefaultAlreadyExists             string
	HostValidationDomainMustBeEmptyForDefault      string
	HostValidationDuplicatedRoutePriority          string
	HostValidationDuplicatedSourcePath             string
	HostValidationDuplicatedVpnName                string
	HostValidationIntegrationRequired              string
	HostValidationJsMainFunctionRequired           string
	HostValidationSourceCodeRequired               string
	HostValidationStaticResponseRequired           string
	HostValidationTargetURIRequired                string
	HostValidationVpnDisabled                      string
	HostValidationVpnNotFound                      string
	IntegrationValidationInUse                     string
	LetsEncryptValidationTosRequired               string
	StreamValidationAtLeastOneBackend              string
	StreamValidationAtLeastOneDomain               string
	StreamValidationFeatureOnlyForTCP              string
	StreamValidationNilStream                      string
	StreamValidationPortNotAllowedForSocket        string
	StreamValidationPortRequired                   string
	StreamValidationRoutesRequiredForSni           string
	UserValidationAtLeastReadOnly                  string
	UserValidationCannotDisableSelf                string
	UserValidationCannotHaveWriteAccess            string
	UserValidationCurrentPasswordMismatch          string
	UserValidationDuplicatedUsername               string
	UserValidationInvalidAccessLevel               string
	VpnValidationInUse                             string
}
