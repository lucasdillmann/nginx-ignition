package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBuildCommands(t *testing.T) {
	t.Run("builds commands with all service methods", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := NewMockRepository(ctrl)
		commands := buildCommands(repo)

		assert.NotNil(t, commands)
		assert.NotNil(t, commands.Get)
		assert.NotNil(t, commands.Save)
		assert.NotNil(t, commands.Delete)
		assert.NotNil(t, commands.Exists)
		assert.NotNil(t, commands.List)
		assert.NotNil(t, commands.GetAvailableDrivers)
		assert.NotNil(t, commands.GetOption)
		assert.NotNil(t, commands.GetOptionURL)
		assert.NotNil(t, commands.ListOptions)
	})
}
