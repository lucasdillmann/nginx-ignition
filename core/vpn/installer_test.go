package vpn

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
)

func Test_BuildCommands(t *testing.T) {
	t.Run("builds commands with all service methods", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cfg := &configuration.Configuration{}
		repo := NewMockRepository(ctrl)
		commands := buildCommands(cfg, repo)

		assert.NotNil(t, commands)
		assert.NotNil(t, commands.Get)
		assert.NotNil(t, commands.Save)
		assert.NotNil(t, commands.Delete)
		assert.NotNil(t, commands.Exists)
		assert.NotNil(t, commands.GetAvailableDrivers)
		assert.NotNil(t, commands.Start)
		assert.NotNil(t, commands.Reload)
		assert.NotNil(t, commands.Stop)
		assert.NotNil(t, commands.List)
	})
}
