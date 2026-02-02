package logline

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/nginx"
)

func TestToResponseDTO(t *testing.T) {
	t.Run("converts LogLine without highlight correctly", func(t *testing.T) {
		input := nginx.LogLine{
			LineNumber: 10,
			Contents:   "some log message",
			Highlight:  nil,
		}

		result := ToResponseDTO(input)

		assert.Equal(t, 10, result.LineNumber)
		assert.Equal(t, "some log message", result.Contents)
		assert.Nil(t, result.Highlight)
	})

	t.Run("converts LogLine with highlight correctly", func(t *testing.T) {
		input := nginx.LogLine{
			LineNumber: 20,
			Contents:   "highlighted message",
			Highlight: &nginx.LogLineHighlight{
				Start: 5,
				End:   15,
			},
		}

		result := ToResponseDTO(input)

		assert.Equal(t, 20, result.LineNumber)
		assert.Equal(t, "highlighted message", result.Contents)
		assert.NotNil(t, result.Highlight)
		assert.Equal(t, 5, result.Highlight.Start)
		assert.Equal(t, 15, result.Highlight.End)
	})
}

func TestToResponseDTOs(t *testing.T) {
	t.Run("converts multiple LogLines correctly", func(t *testing.T) {
		input := []nginx.LogLine{
			{LineNumber: 1, Contents: "line 1"},
			{LineNumber: 2, Contents: "line 2"},
		}

		result := ToResponseDTOs(input)

		assert.Len(t, result, 2)
		assert.Equal(t, 1, result[0].LineNumber)
		assert.Equal(t, "line 1", result[0].Contents)
		assert.Equal(t, 2, result[1].LineNumber)
		assert.Equal(t, "line 2", result[1].Contents)
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		var input []nginx.LogLine
		result := ToResponseDTOs(input)
		assert.Empty(t, result)
	})
}
