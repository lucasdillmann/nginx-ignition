package authorization

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_CurrentSubject(t *testing.T) {
	t.Run("returns the subject set in the context", func(t *testing.T) {
		ctx := &gin.Context{}
		expected := &Subject{TokenID: "token-id"}
		ctx.Set(RequestSubject, expected)

		assert.Same(t, expected, CurrentSubject(ctx))
	})

	t.Run("returns nil when no subject is set", func(t *testing.T) {
		ctx := &gin.Context{}

		assert.Nil(t, CurrentSubject(ctx))
	})
}
