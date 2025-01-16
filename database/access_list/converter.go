package access_list

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"github.com/google/uuid"
)

func toDomain(model *accessListModel) *access_list.AccessList {
	entries := make([]access_list.AccessListEntry, len(model.EntrySets))
	for index, entry := range model.EntrySets {
		entries[index] = access_list.AccessListEntry{
			Priority:      entry.Priority,
			Outcome:       access_list.Outcome(entry.Outcome),
			SourceAddress: entry.SourceAddresses,
		}
	}

	credentials := make([]access_list.AccessListCredentials, len(model.Credentials))
	for index, credential := range model.Credentials {
		credentials[index] = access_list.AccessListCredentials{
			Username: credential.Username,
			Password: credential.Password,
		}
	}

	return &access_list.AccessList{
		ID:                          model.ID,
		Name:                        model.Name,
		Realm:                       model.Realm,
		SatisfyAll:                  model.SatisfyAll,
		DefaultOutcome:              access_list.Outcome(model.DefaultOutcome),
		Entries:                     entries,
		Credentials:                 credentials,
		ForwardAuthenticationHeader: model.ForwardAuthenticationHeader,
	}
}

func toModel(domain *access_list.AccessList) *accessListModel {
	entrySets := make([]*entrySetModel, len(domain.Entries))
	for index, entry := range domain.Entries {
		entrySets[index] = &entrySetModel{
			ID:              uuid.New(),
			AccessListID:    domain.ID,
			Priority:        entry.Priority,
			Outcome:         string(entry.Outcome),
			SourceAddresses: entry.SourceAddress,
		}
	}

	credentials := make([]*credentialsModel, len(domain.Credentials))
	for index, cred := range domain.Credentials {
		credentials[index] = &credentialsModel{
			ID:           uuid.New(),
			AccessListID: domain.ID,
			Username:     cred.Username,
			Password:     cred.Password,
		}
	}

	return &accessListModel{
		ID:                          domain.ID,
		Name:                        domain.Name,
		Realm:                       domain.Realm,
		DefaultOutcome:              string(domain.DefaultOutcome),
		ForwardAuthenticationHeader: domain.ForwardAuthenticationHeader,
		SatisfyAll:                  domain.SatisfyAll,
		Credentials:                 credentials,
		EntrySets:                   entrySets,
	}
}
