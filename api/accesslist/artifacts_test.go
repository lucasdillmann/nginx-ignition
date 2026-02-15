package accesslist

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func newAccessListRequestDTO() accessListRequestDTO {
	return accessListRequestDTO{
		Name:                        new("Test List"),
		Realm:                       new("Test Realm"),
		SatisfyAll:                  new(true),
		DefaultOutcome:              new(accesslist.AllowOutcome),
		ForwardAuthenticationHeader: new(true),
		Entries: []entrySetDTO{
			{
				Priority:        new(1),
				Outcome:         new(accesslist.AllowOutcome),
				SourceAddresses: []string{"192.168.1.1"},
			},
		},
		Credentials: []credentialsDTO{
			{
				Username: new("user1"),
				Password: new("pass1"),
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
