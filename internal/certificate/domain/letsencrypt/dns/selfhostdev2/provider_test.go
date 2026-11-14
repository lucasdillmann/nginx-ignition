package selfhostdev2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseAPIKeys(t *testing.T) {
	t.Run("parses a single pair", func(t *testing.T) {
		credentials, err := parseAPIKeys(context.Background(), "example.com:42.base64")

		require.NoError(t, err)
		assert.Equal(t, map[string]string{"example.com": "42.base64"}, credentials)
	})

	t.Run("parses multiple pairs and trims whitespace", func(t *testing.T) {
		credentials, err := parseAPIKeys(
			context.Background(),
			" example.com : 42.base64 , example.org:99.abc ",
		)

		require.NoError(t, err)
		assert.Equal(t, map[string]string{
			"example.com": "42.base64",
			"example.org": "99.abc",
		}, credentials)
	})

	for name, raw := range map[string]string{
		"missing separator": "example.com",
		"missing domain":    ":42.base64",
		"missing key":       "example.com:",
		"too many parts":    "example.com:42.base64:extra",
		"empty input":       "",
		"empty entry":       "example.com:42.base64,",
	} {
		t.Run("rejects "+name, func(t *testing.T) {
			credentials, err := parseAPIKeys(context.Background(), raw)

			require.Error(t, err)
			assert.Nil(t, credentials)
		})
	}
}

func Test_parseRecordsMapping(t *testing.T) {
	t.Run("parses a single record ID", func(t *testing.T) {
		mapping, err := parseRecordsMapping(context.Background(), "example.com:123")

		require.NoError(t, err)
		require.Len(t, mapping, 1)

		sequence := mapping["example.com"]
		require.NotNil(t, sequence)
		assert.Equal(t, int64(123), sequence.Next())
		assert.Equal(t, int64(123), sequence.Next())
	})

	t.Run("parses multiple record IDs", func(t *testing.T) {
		mapping, err := parseRecordsMapping(context.Background(), "example.com:123:456")

		require.NoError(t, err)
		require.Len(t, mapping, 1)

		sequence := mapping["example.com"]
		require.NotNil(t, sequence)
		assert.Equal(t, int64(123), sequence.Next())
		assert.Equal(t, int64(456), sequence.Next())
		assert.Equal(t, int64(123), sequence.Next())
	})

	t.Run("parses multiple domains and trims whitespace", func(t *testing.T) {
		mapping, err := parseRecordsMapping(
			context.Background(),
			" example.com : 123 : 456 , example.org:789 ",
		)

		require.NoError(t, err)
		require.Len(t, mapping, 2)

		sequence := mapping["example.com"]
		require.NotNil(t, sequence)
		assert.Equal(t, int64(123), sequence.Next())

		sequence = mapping["example.org"]
		require.NotNil(t, sequence)
		assert.Equal(t, int64(789), sequence.Next())
	})

	for name, raw := range map[string]string{
		"missing separator": "example.com",
		"missing domain":    ":123",
		"missing ID":        "example.com:",
		"non-numeric ID":    "example.com:abc",
		"non-integer ID":    "example.com:12.3",
		"too many IDs":      "example.com:123:456:789",
		"empty input":       "",
		"empty entry":       "example.com:123,",
	} {
		t.Run("rejects "+name, func(t *testing.T) {
			mapping, err := parseRecordsMapping(context.Background(), raw)

			require.Error(t, err)
			assert.Nil(t, mapping)
		})
	}
}
