package access_list

import (
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"net"
	"strconv"
	"strings"
)

type validator struct {
	delegate *validation.ConsistencyValidator
}

func newValidator() *validator {
	return &validator{
		delegate: validation.NewValidator(),
	}
}

func (v *validator) validate(accessList *AccessList) error {
	if strings.TrimSpace(accessList.Name) == "" {
		v.delegate.Add("name", validation.ValueMissingMessage)
	}

	var knownUsernames map[string]bool
	for index, value := range accessList.Credentials {
		v.validateCredentials(index, &value, &knownUsernames)
	}

	var knownPriorities map[int]bool
	for index, value := range accessList.Entries {
		v.validateEntry(index, &value, &knownPriorities)
	}

	return v.delegate.Result()
}

func (v *validator) validateEntry(
	index int,
	entry *AccessListEntry,
	knownUsernames *map[int]bool,
) {
	path := "entries[" + strconv.Itoa(index) + "]"
	if (*knownUsernames)[entry.Priority] {
		v.delegate.Add(path+".priority", "Value is duplicated")
	} else {
		(*knownUsernames)[entry.Priority] = true
	}

	if entry.Priority < 0 {
		v.delegate.Add(path+".priority", "Value must be 0 or greater")
	}

	if len(entry.SourceAddress) == 0 {
		v.delegate.Add(path+".sourceAddress", validation.ValueMissingMessage)
	}

	for addressIndex, address := range entry.SourceAddress {
		if singleIp := net.ParseIP(*address); singleIp != nil {
			continue
		}

		if _, _, err := net.ParseCIDR(*address); err == nil {
			continue
		}

		v.delegate.Add(
			path+".sourceAddress["+strconv.Itoa(addressIndex)+"]",
			"Address \""+*address+"\" is not a valid IPv4 or IPv6 address or range",
		)
	}
}

func (v *validator) validateCredentials(
	index int,
	credentials *AccessListCredentials,
	knownUsernames *map[string]bool,
) {
	path := "credentials[" + strconv.Itoa(index) + "].username"

	if strings.TrimSpace(credentials.Username) == "" {
		v.delegate.Add(path, validation.ValueMissingMessage)
	}

	if (*knownUsernames)[credentials.Username] {
		v.delegate.Add(path, credentials.Username)
	} else {
		(*knownUsernames)[credentials.Username] = true
	}
}
