package accesslist

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/accesslist"
)

func toDTO(accessList *accesslist.AccessList) *accessListResponseDTO {
	if accessList == nil {
		return nil
	}

	entries := make([]entrySetDTO, 0)
	for _, entry := range accessList.Entries {
		entries = append(entries, entrySetDTO{
			Priority:        &entry.Priority,
			Outcome:         &entry.Outcome,
			SourceAddresses: entry.SourceAddress,
		})
	}

	credentials := make([]credentialsDTO, 0)
	for _, credential := range accessList.Credentials {
		credentials = append(credentials, credentialsDTO{
			Username: &credential.Username,
			Password: &credential.Password,
		})
	}

	return &accessListResponseDTO{
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

func toDomain(request *accessListRequestDTO) *accesslist.AccessList {
	if request == nil {
		return nil
	}

	entries := make([]accesslist.Entry, 0)
	for _, entry := range request.Entries {
		entries = append(entries, accesslist.Entry{
			Priority:      *entry.Priority,
			Outcome:       *entry.Outcome,
			SourceAddress: entry.SourceAddresses,
		})
	}

	credentials := make([]accesslist.Credentials, 0)
	for _, credential := range request.Credentials {
		credentials = append(credentials, accesslist.Credentials{
			Username: *credential.Username,
			Password: *credential.Password,
		})
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
