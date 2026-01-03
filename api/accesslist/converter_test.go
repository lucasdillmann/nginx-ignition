package accesslist

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/accesslist"
)

func Test_Converter(t *testing.T) {
	t.Run("toDTO", func(t *testing.T) {
		t.Run("returns nil when input is nil", func(t *testing.T) {
			result := toDTO(nil)
			assert.Nil(t, result)
		})

		t.Run("converts full access list to DTO", func(t *testing.T) {
			id := uuid.New()
			input := &accesslist.AccessList{
				ID:                          id,
				Name:                        "Test List",
				Realm:                       "Test Realm",
				SatisfyAll:                  true,
				DefaultOutcome:              accesslist.AllowOutcome,
				ForwardAuthenticationHeader: true,
				Entries: []accesslist.Entry{
					{
						Priority:      1,
						Outcome:       accesslist.AllowOutcome,
						SourceAddress: []string{"192.168.1.1"},
					},
				},
				Credentials: []accesslist.Credentials{
					{
						Username: "user1",
						Password: "pass1",
					},
				},
			}

			result := toDTO(input)

			assert.NotNil(t, result)
			assert.Equal(t, input.ID, result.ID)
			assert.Equal(t, input.Name, result.Name)
			assert.Equal(t, &input.Realm, result.Realm)
			assert.Equal(t, input.SatisfyAll, result.SatisfyAll)
			assert.Equal(t, input.DefaultOutcome, result.DefaultOutcome)
			assert.Equal(t, input.ForwardAuthenticationHeader, result.ForwardAuthenticationHeader)

			assert.Len(t, result.Entries, 1)
			assert.Equal(t, input.Entries[0].Priority, *result.Entries[0].Priority)
			assert.Equal(t, input.Entries[0].Outcome, *result.Entries[0].Outcome)
			assert.Equal(t, input.Entries[0].SourceAddress, result.Entries[0].SourceAddresses)

			assert.Len(t, result.Credentials, 1)
			assert.Equal(t, input.Credentials[0].Username, *result.Credentials[0].Username)
			assert.Equal(t, input.Credentials[0].Password, *result.Credentials[0].Password)
		})
	})

	t.Run("toDomain", func(t *testing.T) {
		t.Run("returns nil when input is nil", func(t *testing.T) {
			result := toDomain(nil)
			assert.Nil(t, result)
		})

		t.Run("converts DTO to domain object", func(t *testing.T) {
			name := "Test List"
			realm := "Test Realm"
			satisfyAll := true
			defaultOutcome := accesslist.AllowOutcome
			forwardAuth := true
			priority := 1
			outcome := accesslist.AllowOutcome
			username := "user1"
			password := "pass1"

			input := &accessListRequestDTO{
				Name:                        &name,
				Realm:                       &realm,
				SatisfyAll:                  &satisfyAll,
				DefaultOutcome:              &defaultOutcome,
				ForwardAuthenticationHeader: &forwardAuth,
				Entries: []entrySetDTO{
					{
						Priority:        &priority,
						Outcome:         &outcome,
						SourceAddresses: []string{"192.168.1.1"},
					},
				},
				Credentials: []credentialsDTO{
					{
						Username: &username,
						Password: &password,
					},
				},
			}

			result := toDomain(input)

			assert.NotNil(t, result)
			assert.NotEqual(t, uuid.Nil, result.ID)
			assert.Equal(t, *input.Name, result.Name)
			assert.Equal(t, *input.Realm, result.Realm)
			assert.Equal(t, *input.SatisfyAll, result.SatisfyAll)
			assert.Equal(t, *input.DefaultOutcome, result.DefaultOutcome)
			assert.Equal(t, *input.ForwardAuthenticationHeader, result.ForwardAuthenticationHeader)

			assert.Len(t, result.Entries, 1)
			assert.Equal(t, *input.Entries[0].Priority, result.Entries[0].Priority)
			assert.Equal(t, *input.Entries[0].Outcome, result.Entries[0].Outcome)
			assert.Equal(t, input.Entries[0].SourceAddresses, result.Entries[0].SourceAddress)

			assert.Len(t, result.Credentials, 1)
			assert.Equal(t, *input.Credentials[0].Username, result.Credentials[0].Username)
			assert.Equal(t, *input.Credentials[0].Password, result.Credentials[0].Password)
		})
	})
}
