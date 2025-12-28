package accesslist

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/accesslist"
)

type accessListRequestDTO struct {
	Name                        *string             `json:"name"`
	Realm                       *string             `json:"realm"`
	SatisfyAll                  *bool               `json:"satisfyAll"`
	DefaultOutcome              *accesslist.Outcome `json:"defaultOutcome"`
	Entries                     []entrySetDTO       `json:"entries"`
	ForwardAuthenticationHeader *bool               `json:"forwardAuthenticationHeader"`
	Credentials                 []credentialsDTO    `json:"credentials"`
}

type accessListResponseDTO struct {
	Realm                       *string            `json:"realm"`
	Name                        string             `json:"name"`
	DefaultOutcome              accesslist.Outcome `json:"defaultOutcome"`
	Entries                     []entrySetDTO      `json:"entries"`
	Credentials                 []credentialsDTO   `json:"credentials"`
	ID                          uuid.UUID          `json:"id"`
	SatisfyAll                  bool               `json:"satisfyAll"`
	ForwardAuthenticationHeader bool               `json:"forwardAuthenticationHeader"`
}

type entrySetDTO struct {
	Priority        *int                `json:"priority"`
	Outcome         *accesslist.Outcome `json:"outcome"`
	SourceAddresses []string            `json:"sourceAddresses"`
}

type credentialsDTO struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}
