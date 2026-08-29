package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Std(t *testing.T) {
	t.Run("returns standard logger", func(t *testing.T) {
		stdLogger := Std()
		assert.NotNil(t, stdLogger)
	})
}
