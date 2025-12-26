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
	Name                        string
	Realm                       string
	DefaultOutcome              Outcome
	Entries                     []Entry
	Credentials                 []Credentials
	ID                          uuid.UUID
	SatisfyAll                  bool
	ForwardAuthenticationHeader bool
}

type Entry struct {
	Outcome       Outcome
	SourceAddress []string
	Priority      int
}

type Credentials struct {
	Username string
	Password string
}
