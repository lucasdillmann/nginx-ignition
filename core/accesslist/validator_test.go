package accesslist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validator(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		t.Run("valid access list passes", func(t *testing.T) {
			accessList := newAccessList()
			accessListValidator := newValidator()

			err := accessListValidator.validate(accessList)

			assert.NoError(t, err)
		})

		t.Run("empty name fails", func(t *testing.T) {
			accessList := newAccessList()
			accessList.Name = ""
			accessListValidator := newValidator()

			err := accessListValidator.validate(accessList)

			assert.Error(t, err)
		})

		t.Run("whitespace-only name fails", func(t *testing.T) {
			accessList := newAccessList()
			accessList.Name = "   "
			accessListValidator := newValidator()

			err := accessListValidator.validate(accessList)

			assert.Error(t, err)
		})
	})

	t.Run("validateEntry", func(t *testing.T) {
		t.Run("valid entry passes", func(t *testing.T) {
			entry := newEntry()
			knownPriorities := map[int]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateEntry(0, entry, &knownPriorities)

			assert.NoError(t, accessListValidator.delegate.Result())
		})

		t.Run("duplicate priority fails", func(t *testing.T) {
			entry1 := newEntry()
			entry2 := newEntry()
			entry2.Priority = entry1.Priority
			knownPriorities := map[int]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateEntry(0, entry1, &knownPriorities)
			accessListValidator.validateEntry(1, entry2, &knownPriorities)

			assert.Error(t, accessListValidator.delegate.Result())
		})

		t.Run("negative priority fails", func(t *testing.T) {
			entry := newEntry()
			entry.Priority = -1
			knownPriorities := map[int]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateEntry(0, entry, &knownPriorities)

			assert.Error(t, accessListValidator.delegate.Result())
		})

		t.Run("zero priority passes", func(t *testing.T) {
			entry := newEntry()
			entry.Priority = 0
			knownPriorities := map[int]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateEntry(0, entry, &knownPriorities)

			assert.NoError(t, accessListValidator.delegate.Result())
		})

		t.Run("empty source address fails", func(t *testing.T) {
			entry := newEntry()
			entry.SourceAddress = []string{}
			knownPriorities := map[int]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateEntry(0, entry, &knownPriorities)

			assert.Error(t, accessListValidator.delegate.Result())
		})

		t.Run("valid IPv4 address passes", func(t *testing.T) {
			entry := newEntry()
			entry.SourceAddress = []string{"192.168.1.1"}
			knownPriorities := map[int]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateEntry(0, entry, &knownPriorities)

			assert.NoError(t, accessListValidator.delegate.Result())
		})

		t.Run("valid IPv6 address passes", func(t *testing.T) {
			entry := newEntry()
			entry.SourceAddress = []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"}
			knownPriorities := map[int]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateEntry(0, entry, &knownPriorities)

			assert.NoError(t, accessListValidator.delegate.Result())
		})

		t.Run("valid CIDR range passes", func(t *testing.T) {
			entry := newEntry()
			entry.SourceAddress = []string{"192.168.1.0/24"}
			knownPriorities := map[int]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateEntry(0, entry, &knownPriorities)

			assert.NoError(t, accessListValidator.delegate.Result())
		})

		t.Run("invalid address fails", func(t *testing.T) {
			entry := newEntry()
			entry.SourceAddress = []string{"invalid.address"}
			knownPriorities := map[int]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateEntry(0, entry, &knownPriorities)

			assert.Error(t, accessListValidator.delegate.Result())
		})

		t.Run("multiple addresses validates all", func(t *testing.T) {
			entry := newEntry()
			entry.SourceAddress = []string{"192.168.1.1", "10.0.0.0/8", "invalid"}
			knownPriorities := map[int]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateEntry(0, entry, &knownPriorities)

			assert.Error(t, accessListValidator.delegate.Result())
		})
	})

	t.Run("validateCredentials", func(t *testing.T) {
		t.Run("valid credentials pass", func(t *testing.T) {
			credentials := newCredentials()
			knownUsernames := map[string]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateCredentials(0, credentials, &knownUsernames)

			assert.NoError(t, accessListValidator.delegate.Result())
		})

		t.Run("empty username fails", func(t *testing.T) {
			credentials := newCredentials()
			credentials.Username = ""
			knownUsernames := map[string]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateCredentials(0, credentials, &knownUsernames)

			assert.Error(t, accessListValidator.delegate.Result())
		})

		t.Run("whitespace-only username fails", func(t *testing.T) {
			credentials := newCredentials()
			credentials.Username = "   "
			knownUsernames := map[string]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateCredentials(0, credentials, &knownUsernames)

			assert.Error(t, accessListValidator.delegate.Result())
		})

		t.Run("duplicate username fails", func(t *testing.T) {
			credentials1 := newCredentials()
			credentials2 := newCredentials()
			knownUsernames := map[string]bool{}
			accessListValidator := newValidator()

			accessListValidator.validateCredentials(0, credentials1, &knownUsernames)
			accessListValidator.validateCredentials(1, credentials2, &knownUsernames)

			assert.Error(t, accessListValidator.delegate.Result())
		})
	})
}
