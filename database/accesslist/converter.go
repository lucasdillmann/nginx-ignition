package accesslist

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/common/pointers"
)

func toDomain(model *accessListModel) *accesslist.AccessList {
	entries := make([]accesslist.AccessListEntry, len(model.EntrySets))
	for index, entry := range model.EntrySets {
		entries[index] = accesslist.AccessListEntry{
			Priority:      entry.Priority,
			Outcome:       accesslist.Outcome(entry.Outcome),
			SourceAddress: pointers.Reference(entry.SourceAddresses),
		}
	}

	credentials := make([]accesslist.AccessListCredentials, len(model.Credentials))
	for index, credential := range model.Credentials {
		credentials[index] = accesslist.AccessListCredentials{
			Username: credential.Username,
			Password: credential.Password,
		}
	}

	return &accesslist.AccessList{
		ID:                          model.ID,
		Name:                        model.Name,
		Realm:                       model.Realm,
		SatisfyAll:                  model.SatisfyAll,
		DefaultOutcome:              accesslist.Outcome(model.DefaultOutcome),
		Entries:                     entries,
		Credentials:                 credentials,
		ForwardAuthenticationHeader: model.ForwardAuthenticationHeader,
	}
}

func toModel(domain *accesslist.AccessList) *accessListModel {
	entrySets := make([]*entrySetModel, len(domain.Entries))
	for index, entry := range domain.Entries {
		entrySets[index] = &entrySetModel{
			ID:              uuid.New(),
			AccessListID:    domain.ID,
			Priority:        entry.Priority,
			Outcome:         string(entry.Outcome),
			SourceAddresses: pointers.Dereference(entry.SourceAddress),
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
