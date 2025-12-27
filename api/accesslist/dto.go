package accesslist

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/accesslist"
)

type accessListRequestDto struct {
	Name                        *string             `json:"name"`
	Realm                       *string             `json:"realm"`
	SatisfyAll                  *bool               `json:"satisfyAll"`
	DefaultOutcome              *accesslist.Outcome `json:"defaultOutcome"`
	Entries                     []entrySetDto       `json:"entries"`
	ForwardAuthenticationHeader *bool               `json:"forwardAuthenticationHeader"`
	Credentials                 []credentialsDto    `json:"credentials"`
}

type accessListResponseDto struct {
	Realm                       *string            `json:"realm"`
	Name                        string             `json:"name"`
	DefaultOutcome              accesslist.Outcome `json:"defaultOutcome"`
	Entries                     []entrySetDto      `json:"entries"`
	Credentials                 []credentialsDto   `json:"credentials"`
	ID                          uuid.UUID          `json:"id"`
	SatisfyAll                  bool               `json:"satisfyAll"`
	ForwardAuthenticationHeader bool               `json:"forwardAuthenticationHeader"`
}

type entrySetDto struct {
	Priority        *int                `json:"priority"`
	Outcome         *accesslist.Outcome `json:"outcome"`
	SourceAddresses []string            `json:"sourceAddresses"`
}

type credentialsDto struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}
