package accesslist

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func newAccessListRequestDTO() accessListRequestDTO {
	name := "Test List"
	realm := "Test Realm"
	satisfyAll := true
	defaultOutcome := accesslist.AllowOutcome
	forwardAuth := true
	priority := 1
	outcome := accesslist.AllowOutcome
	username := "user1"
	password := "pass1"

	return accessListRequestDTO{
		Name:                        &name,
		Realm:                       &realm,
		SatisfyAll:                  &satisfyAll,
		DefaultOutcome:              &defaultOutcome,
		ForwardAuthenticationHeader: &forwardAuth,
		Entries: []entrySetDTO{
			{
				Priority:        &priority,
				Outcome:         &outcome,
				SourceAddresses: []string{"192.168.1.1"},
			},
		},
		Credentials: []credentialsDTO{
			{
				Username: &username,
				Password: &password,
			},
		},
	}
}

func newAccessList() *accesslist.AccessList {
	return &accesslist.AccessList{
		ID:                          uuid.New(),
		Name:                        "Test List",
		Realm:                       "Test Realm",
		SatisfyAll:                  true,
		DefaultOutcome:              accesslist.AllowOutcome,
		ForwardAuthenticationHeader: true,
		Entries: []accesslist.Entry{
			{
				Priority:      1,
				Outcome:       accesslist.AllowOutcome,
				SourceAddress: []string{"192.168.1.1"},
			},
		},
		Credentials: []accesslist.Credentials{
			{
				Username: "user1",
				Password: "pass1",
			},
		},
	}
}

func newAccessListPage() *pagination.Page[accesslist.AccessList] {
	return pagination.Of([]accesslist.AccessList{
		{
			Name: "Test",
		},
	})
}
