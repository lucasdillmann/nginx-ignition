package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
)

func Test_BuildCommands(t *testing.T) {
	t.Run("builds commands and service with all methods", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := NewMockRepository(ctrl)
		cfg := &configuration.Configuration{}
		svc, commands := buildCommands(repo, cfg)

		assert.NotNil(t, svc)
		assert.NotNil(t, commands)
		assert.NotNil(t, commands.Authenticate)
		assert.NotNil(t, commands.Delete)
		assert.NotNil(t, commands.Get)
		assert.NotNil(t, commands.GetCount)
		assert.NotNil(t, commands.GetStatus)
		assert.NotNil(t, commands.List)
		assert.NotNil(t, commands.Save)
		assert.NotNil(t, commands.UpdatePassword)
		assert.NotNil(t, commands.OnboardingCompleted)
	})
}
