package container

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Init(t *testing.T) {
	t.Run("initializes container", func(t *testing.T) {
		Init(t.Context())

		assert.NotNil(t, delegate)
	})
}

func Test_Provide(t *testing.T) {
	t.Run("provides function", func(t *testing.T) {
		Init(t.Context())

		provider := func() string {
			return "test"
		}

		err := Provide(provider)

		assert.NoError(t, err)
	})

	t.Run("returns error for invalid provider", func(t *testing.T) {
		Init(t.Context())

		err := Provide("not-a-function")

		assert.Error(t, err)
	})
}

func Test_Singleton(t *testing.T) {
	t.Run("provides singleton value", func(t *testing.T) {
		Init(t.Context())

		value := "test-singleton"
		err := Singleton(value)

		assert.NoError(t, err)

		result := Get[string]()
		assert.Equal(t, value, result)
	})
}

func Test_Run(t *testing.T) {
	t.Run("invokes function", func(t *testing.T) {
		Init(t.Context())

		called := false
		fn := func() {
			called = true
		}

		err := Run(fn)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("returns error when function fails", func(t *testing.T) {
		Init(t.Context())

		fn := func() error {
			return errors.New("invocation error")
		}

		err := Run(fn)

		assert.Error(t, err)
	})
}

func Test_Get(t *testing.T) {
	t.Run("retrieves provided value", func(t *testing.T) {
		Init(t.Context())

		expectedValue := "test-value"
		_ = Provide(func() string {
			return expectedValue
		})

		result := Get[string]()

		assert.Equal(t, expectedValue, result)
	})
}
