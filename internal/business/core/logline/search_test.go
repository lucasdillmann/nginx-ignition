package logline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Search(t *testing.T) {
	t.Run("returns all lines when query is empty", func(t *testing.T) {
		lines := []LogLine{
			{Contents: "line 1", LineNumber: 1},
			{Contents: "line 2", LineNumber: 2},
		}

		result, err := Search(lines, "", 0)

		assert.NoError(t, err)
		assert.Equal(t, lines, result)
	})

	t.Run("returns empty slice when no match is found", func(t *testing.T) {
		lines := []LogLine{
			{Contents: "line 1", LineNumber: 1},
		}

		result, err := Search(lines, "non-existent", 0)

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("performs case-insensitive search", func(t *testing.T) {
		lines := []LogLine{
			{Contents: "Some LOG Content", LineNumber: 1},
		}

		result, err := Search(lines, "log", 0)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Some LOG Content", result[0].Contents)
	})

	t.Run("replaces spaces with catch-all statement", func(t *testing.T) {
		lines := []LogLine{
			{Contents: "starting middleware ending", LineNumber: 1},
		}

		result, err := Search(lines, "starting ending", 0)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("fills highlight attribute when a match is found", func(t *testing.T) {
		lines := []LogLine{
			{Contents: "prefix MATCH suffix", LineNumber: 1},
		}

		result, err := Search(lines, "match", 0)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.NotNil(t, result[0].Highlight)
		assert.Equal(t, 7, result[0].Highlight.Start)
		assert.Equal(t, 12, result[0].Highlight.End)
	})

	t.Run("does not fill highlight for surrounding lines that do not match", func(t *testing.T) {
		lines := []LogLine{
			{Contents: "context", LineNumber: 1},
			{Contents: "match", LineNumber: 2},
		}

		result, err := Search(lines, "match", 1)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Nil(t, result[0].Highlight)
		assert.NotNil(t, result[1].Highlight)
	})

	t.Run("surroundingLines", func(t *testing.T) {
		lines := []LogLine{
			{Contents: "line 1", LineNumber: 1},
			{Contents: "line 2", LineNumber: 2},
			{Contents: "match", LineNumber: 3},
			{Contents: "line 4", LineNumber: 4},
			{Contents: "line 5", LineNumber: 5},
		}

		t.Run("includes specified number of surrounding lines", func(t *testing.T) {
			result, err := Search(lines, "match", 1)

			assert.NoError(t, err)
			assert.Len(t, result, 3)
			assert.Equal(t, 2, result[0].LineNumber)
			assert.Equal(t, 3, result[1].LineNumber)
			assert.Equal(t, 4, result[2].LineNumber)
		})

		t.Run("limits surrounding lines to maximum of 10", func(t *testing.T) {
			manyLines := make([]LogLine, 25)
			for i := range manyLines {
				manyLines[i] = LogLine{Contents: "line", LineNumber: i + 1}
			}
			manyLines[12].Contents = "match"

			result, err := Search(manyLines, "match", 50)

			assert.NoError(t, err)
			assert.Len(t, result, 21)
			assert.Equal(t, 3, result[0].LineNumber)
			assert.Equal(t, 23, result[len(result)-1].LineNumber)
		})

		t.Run("merges overlapping surrounding ranges", func(t *testing.T) {
			lines := []LogLine{
				{Contents: "match 1", LineNumber: 1},
				{Contents: "line 2", LineNumber: 2},
				{Contents: "match 2", LineNumber: 3},
			}

			result, err := Search(lines, "match", 1)

			assert.NoError(t, err)
			assert.Len(t, result, 3)
		})
	})
}
