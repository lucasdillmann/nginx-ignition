package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Init(t *testing.T) {
	t.Run("initializes logger", func(t *testing.T) {
		err := Init()

		assert.NoError(t, err)
		assert.NotNil(t, delegate)
	})
}

func Test_Std(t *testing.T) {
	t.Run("returns standard logger", func(t *testing.T) {
		_ = Init()

		stdLogger := Std()

		assert.NotNil(t, stdLogger)
	})
}
