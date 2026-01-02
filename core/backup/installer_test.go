package backup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_BuildCommands(t *testing.T) {
	t.Run("builds commands with all service methods", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := NewMockRepository(ctrl)
		commands := buildCommands(repo)

		assert.NotNil(t, commands)
		assert.NotNil(t, commands.Get)
	})
}
