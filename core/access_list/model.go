package access_list

import (
	"github.com/google/uuid"
)

type Outcome string

const (
	AllowOutcome Outcome = "ALLOW"
	DenyOutcome  Outcome = "DENY"
)

type AccessList struct {
	ID                          uuid.UUID
	Name                        string
	Realm                       string
	SatisfyAll                  bool
	DefaultOutcome              Outcome
	Entries                     []AccessListEntry
	Credentials                 []AccessListCredentials
	ForwardAuthenticationHeader bool
}

type AccessListEntry struct {
	Priority      int
	Outcome       Outcome
	SourceAddress []*string
}

type AccessListCredentials struct {
	Username string
	Password string
}
