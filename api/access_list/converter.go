package access_list

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"github.com/google/uuid"
)

func toDto(accessList *access_list.AccessList) *accessListResponseDto {
	if accessList == nil {
		return nil
	}

	var entries []entrySetDto
	for _, entry := range accessList.Entries {
		entries = append(entries, entrySetDto{
			Priority:        &entry.Priority,
			Outcome:         &entry.Outcome,
			SourceAddresses: &entry.SourceAddress,
		})
	}

	var credentials []credentialsDto
	for _, credential := range accessList.Credentials {
		credentials = append(credentials, credentialsDto{
			Username: &credential.Username,
			Password: &credential.Password,
		})
	}

	return &accessListResponseDto{
		ID:                          accessList.ID,
		Name:                        accessList.Name,
		Realm:                       &accessList.Realm,
		SatisfyAll:                  accessList.SatisfyAll,
		DefaultOutcome:              accessList.DefaultOutcome,
		Entries:                     entries,
		ForwardAuthenticationHeader: accessList.ForwardAuthenticationHeader,
		Credentials:                 credentials,
	}
}

func toDomain(request *accessListRequestDto) *access_list.AccessList {
	if request == nil {
		return nil
	}

	var entries []access_list.AccessListEntry
	if request.Entries != nil {
		for _, entry := range *request.Entries {
			entries = append(entries, access_list.AccessListEntry{
				Priority:      *entry.Priority,
				Outcome:       *entry.Outcome,
				SourceAddress: *entry.SourceAddresses,
			})
		}
	}

	var credentials []access_list.AccessListCredentials
	if request.Credentials != nil {
		for _, credential := range *request.Credentials {
			credentials = append(credentials, access_list.AccessListCredentials{
				Username: *credential.Username,
				Password: *credential.Password,
			})
		}
	}

	return &access_list.AccessList{
		ID:                          uuid.New(),
		Name:                        *request.Name,
		Realm:                       *request.Realm,
		SatisfyAll:                  *request.SatisfyAll,
		DefaultOutcome:              *request.DefaultOutcome,
		Entries:                     entries,
		ForwardAuthenticationHeader: *request.ForwardAuthenticationHeader,
		Credentials:                 credentials,
	}
}
