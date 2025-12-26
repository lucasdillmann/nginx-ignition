package accesslist

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
	Entries                     []Entry
	Credentials                 []Credentials
	ForwardAuthenticationHeader bool
}

type Entry struct {
	Priority      int
	Outcome       Outcome
	SourceAddress []string
}

type Credentials struct {
	Username string
	Password string
}
