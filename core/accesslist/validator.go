package accesslist

import (
	"context"
	"fmt"
	"net"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type validator struct {
	delegate *validation.ConsistencyValidator
}

func newValidator() *validator {
	return &validator{
		delegate: validation.NewValidator(),
	}
}

func (v *validator) validate(ctx context.Context, accessList *AccessList) error {
	if strings.TrimSpace(accessList.Name) == "" {
		v.delegate.Add("name", i18n.M(ctx, i18n.K.CommonValidationValueMissing))
	}

	knownUsernames := map[string]bool{}
	for index, value := range accessList.Credentials {
		v.validateCredentials(ctx, index, &value, &knownUsernames)
	}

	knownPriorities := map[int]bool{}
	for index, value := range accessList.Entries {
		v.validateEntry(ctx, index, &value, &knownPriorities)
	}

	return v.delegate.Result()
}

func (v *validator) validateEntry(
	ctx context.Context,
	index int,
	entry *Entry,
	knownUsernames *map[int]bool,
) {
	path := fmt.Sprintf("entries[%d]", index)
	if (*knownUsernames)[entry.Priority] {
		v.delegate.Add(path+".priority", i18n.M(ctx, i18n.K.CommonValidationDuplicatedValue))
	} else {
		(*knownUsernames)[entry.Priority] = true
	}

	if entry.Priority < 0 {
		v.delegate.Add(path+".priority", i18n.M(ctx, i18n.K.CommonValidationCannotBeZero))
	}

	if len(entry.SourceAddress) == 0 {
		v.delegate.Add(path+".sourceAddress", i18n.M(ctx, i18n.K.CommonValidationValueMissing))
	}

	for addressIndex, address := range entry.SourceAddress {
		if singleIP := net.ParseIP(address); singleIP != nil {
			continue
		}

		if _, _, err := net.ParseCIDR(address); err == nil {
			continue
		}

		v.delegate.Add(
			fmt.Sprintf("%s.sourceAddress[%d]", path, addressIndex),
			i18n.M(ctx, i18n.K.AccessListValidationInvalidAddress).V("address", address),
		)
	}
}

func (v *validator) validateCredentials(
	ctx context.Context,
	index int,
	credentials *Credentials,
	knownUsernames *map[string]bool,
) {
	path := fmt.Sprintf("credentials[%d].username", index)

	if strings.TrimSpace(credentials.Username) == "" {
		v.delegate.Add(path, i18n.M(ctx, i18n.K.CommonValidationValueMissing))
	}

	if (*knownUsernames)[credentials.Username] {
		v.delegate.Add(path, i18n.Raw(credentials.Username))
	} else {
		(*knownUsernames)[credentials.Username] = true
	}
}
