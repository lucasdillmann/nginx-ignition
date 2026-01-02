package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
)

func Test_BuildCommands(t *testing.T) {
	t.Run("builds commands with all service methods", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := NewMockRepository(ctrl)
		bindingCommands := &binding.Commands{}
		sched := &scheduler.Scheduler{}
		commands := buildCommands(repo, bindingCommands, sched)

		assert.NotNil(t, commands)
		assert.NotNil(t, commands.Get)
		assert.NotNil(t, commands.Save)
	})
}
