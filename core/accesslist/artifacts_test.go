package accesslist

import (
	"github.com/google/uuid"
)

func newAccessList() *AccessList {
	return &AccessList{
		ID:             uuid.New(),
		Name:           "Default Access List",
		DefaultOutcome: AllowOutcome,
		Entries: []Entry{
			{
				Outcome:       DenyOutcome,
				SourceAddress: []string{"192.168.1.50", "10.0.0.0/24"},
				Priority:      10,
			},
			{
				Outcome:       AllowOutcome,
				SourceAddress: []string{"127.0.0.1"},
				Priority:      100,
			},
		},
		Credentials: []Credentials{
			{
				Username: "user1",
				Password: "pass1",
			},
		},
	}
}

func newEntry() *Entry {
	return &Entry{
		Outcome:       AllowOutcome,
		SourceAddress: []string{"192.168.1.1"},
		Priority:      1,
	}
}

func newCredentials() *Credentials {
	return &Credentials{
		Username: "user1",
		Password: "pass1",
	}
}
