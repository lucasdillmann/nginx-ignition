package letsencrypt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_providers(t *testing.T) {
	t.Run("matches the number of provider packages on disk", func(t *testing.T) {
		entries, err := os.ReadDir("./dns")
		require.NoError(t, err)

		packageCount := 0

		for _, entry := range entries {
			if entry.IsDir() {
				packageCount++
			}
		}

		assert.Equal(
			t,
			packageCount,
			len(providers),
			"every package under ./dns must be registered in the providers list",
		)
	})

	t.Run("has unique IDs", func(t *testing.T) {
		seen := make(map[string]int, len(providers))

		for index, provider := range providers {
			if previous, duplicated := seen[provider.ID()]; duplicated {
				t.Errorf(
					"duplicated provider ID %q at indexes %d and %d",
					provider.ID(),
					previous,
					index,
				)

				continue
			}

			seen[provider.ID()] = index
		}
	})
}
