package access_list

import "github.com/google/uuid"

type AccessListOutcome string

const (
	AllowAccessOutcome = AccessListOutcome("ALLOW")
	DenyAccessOutcome  = AccessListOutcome("DENY")
)

type AccessList struct {
	Id                          uuid.UUID
	Name                        string
	Realm                       string
	SatisfyAll                  bool
	DefaultOutcome              AccessListOutcome
	Entries                     []AccessListEntry
	Credentials                 []AccessListCredentials
	ForwardAuthenticationHeader bool
}

type AccessListEntry struct {
	Priority      int64
	Outcome       AccessListOutcome
	SourceAddress []string
}

type AccessListCredentials struct {
	Username string
	Password string
}
