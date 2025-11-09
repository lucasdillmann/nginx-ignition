package access_list

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/access_list"
)

type accessListRequestDto struct {
	Name                        *string              `json:"name"`
	Realm                       *string              `json:"realm"`
	SatisfyAll                  *bool                `json:"satisfyAll"`
	DefaultOutcome              *access_list.Outcome `json:"defaultOutcome"`
	Entries                     []*entrySetDto       `json:"entries"`
	ForwardAuthenticationHeader *bool                `json:"forwardAuthenticationHeader"`
	Credentials                 []*credentialsDto    `json:"credentials"`
}

type accessListResponseDto struct {
	ID                          uuid.UUID           `json:"id"`
	Name                        string              `json:"name"`
	Realm                       *string             `json:"realm"`
	SatisfyAll                  bool                `json:"satisfyAll"`
	DefaultOutcome              access_list.Outcome `json:"defaultOutcome"`
	Entries                     []entrySetDto       `json:"entries"`
	ForwardAuthenticationHeader bool                `json:"forwardAuthenticationHeader"`
	Credentials                 []credentialsDto    `json:"credentials"`
}

type entrySetDto struct {
	Priority        *int                 `json:"priority"`
	Outcome         *access_list.Outcome `json:"outcome"`
	SourceAddresses []*string            `json:"sourceAddresses"`
}

type credentialsDto struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}
