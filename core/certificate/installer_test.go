package certificate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBuildCommands(t *testing.T) {
	t.Run("builds commands and service with all methods", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := NewMockRepository(ctrl)
		commands, svc := buildCommands(repo)

		assert.NotNil(t, commands)
		assert.NotNil(t, svc)
		assert.NotNil(t, commands.Delete)
		assert.NotNil(t, commands.AvailableProviders)
		assert.NotNil(t, commands.Exists)
		assert.NotNil(t, commands.Get)
		assert.NotNil(t, commands.List)
		assert.NotNil(t, commands.Issue)
		assert.NotNil(t, commands.Renew)
	})
}
