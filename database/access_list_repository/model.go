package access_list_repository

import (
	"github.com/google/uuid"
)

type accessListModel struct {
	ID                          uuid.UUID          `bun:"id,pk"`
	Name                        string             `bun:"name,unique,notnull"`
	Realm                       string             `bun:"realm"`
	DefaultOutcome              string             `bun:"default_outcome,notnull"`
	ForwardAuthenticationHeader bool               `bun:"forward_authentication_header,notnull"`
	SatisfyAll                  bool               `bun:"satisfy_all,notnull"`
	Credentials                 []credentialsModel `bun:"rel:has-many,join:id=access_list_id"`
	EntrySets                   []entrySetModel    `bun:"rel:has-many,join:id=access_list_id"`
}

type credentialsModel struct {
	ID           uuid.UUID `bun:"id,pk"`
	AccessListID uuid.UUID `bun:"access_list_id,notnull"`
	Username     string    `bun:"username,notnull"`
	Password     string    `bun:"password,notnull"`
}

type entrySetModel struct {
	ID              uuid.UUID `bun:"id,pk"`
	AccessListID    uuid.UUID `bun:"access_list_id,notnull"`
	Priority        int       `bun:"priority,notnull"`
	Outcome         string    `bun:"outcome,notnull"`
	SourceAddresses []string  `bun:"source_addresses,array,notnull"`
}
