package cfgfiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_flag(t *testing.T) {
	t.Run("returns trueValue when enabled is true", func(t *testing.T) {
		assert.Equal(t, "on", flag(true, "on", "off"))
	})

	t.Run("returns falseValue when enabled is false", func(t *testing.T) {
		assert.Equal(t, "off", flag(false, "on", "off"))
	})
}

func Test_statusFlag(t *testing.T) {
	t.Run("returns on when true", func(t *testing.T) {
		assert.Equal(t, "on", statusFlag(true))
	})

	t.Run("returns off when false", func(t *testing.T) {
		assert.Equal(t, "off", statusFlag(false))
	})
}
