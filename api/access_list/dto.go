package access_list

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"github.com/google/uuid"
)

type accessListRequestDto struct {
	Name                        *string              `json:"name" validate:"required"`
	Realm                       *string              `json:"realm"`
	SatisfyAll                  *bool                `json:"satisfyAll" validate:"required"`
	DefaultOutcome              *access_list.Outcome `json:"defaultOutcome" validate:"required"`
	Entries                     *[]entrySetDto       `json:"entries" validate:"required"`
	ForwardAuthenticationHeader *bool                `json:"forwardAuthenticationHeader" validate:"required"`
	Credentials                 *[]credentialsDto    `json:"credentials" validate:"required"`
}

type accessListResponseDto struct {
	ID                          uuid.UUID           `json:"id" validate:"required"`
	Name                        string              `json:"name" validate:"required"`
	Realm                       *string             `json:"realm"`
	SatisfyAll                  bool                `json:"satisfyAll" validate:"required"`
	DefaultOutcome              access_list.Outcome `json:"defaultOutcome" validate:"required"`
	Entries                     []entrySetDto       `json:"entries" validate:"required"`
	ForwardAuthenticationHeader bool                `json:"forwardAuthenticationHeader" validate:"required"`
	Credentials                 []credentialsDto    `json:"credentials" validate:"required"`
}

type entrySetDto struct {
	Priority        *int                 `json:"priority" validate:"required"`
	Outcome         *access_list.Outcome `json:"outcome" validate:"required"`
	SourceAddresses *[]string            `json:"sourceAddresses" validate:"required"`
}

type credentialsDto struct {
	Username *string `json:"username" validate:"required"`
	Password *string `json:"password" validate:"required"`
}
