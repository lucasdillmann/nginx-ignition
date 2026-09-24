package accesslist

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ToDTO(t *testing.T) {
	t.Run("returns nil when input is nil", func(t *testing.T) {
		response := toDTO(nil)
		assert.Nil(t, response)
	})

	t.Run("converts full access list to DTO", func(t *testing.T) {
		accessList := newAccessList()
		response := toDTO(accessList)

		assert.NotNil(t, response)
		assert.Equal(t, accessList.ID, response.ID)
		assert.Equal(t, accessList.Name, response.Name)
		assert.Equal(t, &accessList.Realm, response.Realm)
		assert.Equal(t, accessList.SatisfyAll, response.SatisfyAll)
		assert.Equal(t, accessList.DefaultOutcome, response.DefaultOutcome)
		assert.Equal(
			t,
			accessList.ForwardAuthenticationHeader,
			response.ForwardAuthenticationHeader,
		)

		assert.Len(t, response.Entries, 1)
		assert.Equal(t, accessList.Entries[0].Priority, *response.Entries[0].Priority)
		assert.Equal(t, accessList.Entries[0].Outcome, *response.Entries[0].Outcome)
		assert.Equal(t, accessList.Entries[0].SourceAddress, response.Entries[0].SourceAddresses)

		assert.Len(t, response.Credentials, 1)
		assert.Equal(t, accessList.Credentials[0].Username, *response.Credentials[0].Username)
		assert.Equal(t, accessList.Credentials[0].Password, *response.Credentials[0].Password)
	})
}

func Test_ToDomain(t *testing.T) {
	t.Run("returns nil when input is nil", func(t *testing.T) {
		accessList := toDomain(nil)
		assert.Nil(t, accessList)
	})

	t.Run("converts DTO to domain object", func(t *testing.T) {
		payload := newAccessListRequestDTO()
		accessList := toDomain(&payload)

		assert.NotNil(t, accessList)
		assert.NotEqual(t, uuid.Nil, accessList.ID)
		assert.Equal(t, *payload.Name, accessList.Name)
		assert.Equal(t, *payload.Realm, accessList.Realm)
		assert.Equal(t, *payload.SatisfyAll, accessList.SatisfyAll)
		assert.Equal(t, *payload.DefaultOutcome, accessList.DefaultOutcome)
		assert.Equal(
			t,
			*payload.ForwardAuthenticationHeader,
			accessList.ForwardAuthenticationHeader,
		)

		assert.Len(t, accessList.Entries, 1)
		assert.Equal(t, *payload.Entries[0].Priority, accessList.Entries[0].Priority)
		assert.Equal(t, *payload.Entries[0].Outcome, accessList.Entries[0].Outcome)
		assert.Equal(t, payload.Entries[0].SourceAddresses, accessList.Entries[0].SourceAddress)

		assert.Len(t, accessList.Credentials, 1)
		assert.Equal(t, *payload.Credentials[0].Username, accessList.Credentials[0].Username)
		assert.Equal(t, *payload.Credentials[0].Password, accessList.Credentials[0].Password)
	})
}
