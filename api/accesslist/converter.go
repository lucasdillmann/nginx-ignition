package accesslist

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/accesslist"
)

func toDto(accessList *accesslist.AccessList) *accessListResponseDto {
	if accessList == nil {
		return nil
	}

	var entries []entrySetDto
	for _, entry := range accessList.Entries {
		entries = append(entries, entrySetDto{
			Priority:        &entry.Priority,
			Outcome:         &entry.Outcome,
			SourceAddresses: entry.SourceAddress,
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

func toDomain(request *accessListRequestDto) *accesslist.AccessList {
	if request == nil {
		return nil
	}

	var entries []accesslist.Entry
	if request.Entries != nil {
		for _, entry := range request.Entries {
			entries = append(entries, accesslist.Entry{
				Priority:      *entry.Priority,
				Outcome:       *entry.Outcome,
				SourceAddress: entry.SourceAddresses,
			})
		}
	}

	var credentials []accesslist.Credentials
	if request.Credentials != nil {
		for _, credential := range request.Credentials {
			credentials = append(credentials, accesslist.Credentials{
				Username: *credential.Username,
				Password: *credential.Password,
			})
		}
	}

	return &accesslist.AccessList{
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
