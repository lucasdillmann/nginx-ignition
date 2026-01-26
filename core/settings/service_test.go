package settings

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
)

func Test_service(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("returns settings when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := &Settings{}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Get(t.Context()).Return(expected, nil)

			bindingCommands := binding.NewMockedCommands(ctrl)
			sched := &scheduler.Scheduler{}

			settingsService := newCommands(repo, bindingCommands, sched)
			result, err := settingsService.Get(t.Context())

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedErr := errors.New("repository error")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Get(t.Context()).Return(nil, expectedErr)

			bindingCommands := binding.NewMockedCommands(ctrl)
			sched := &scheduler.Scheduler{}

			settingsService := newCommands(repo, bindingCommands, sched)
			result, err := settingsService.Get(t.Context())

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("invalid settings returns validation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := newSettings()
			s.Nginx.DefaultContentType = "" // Invalid

			repo := NewMockedRepository(ctrl)
			bindingCommands := binding.NewMockedCommands(ctrl)
			sched := &scheduler.Scheduler{}

			settingsService := newCommands(repo, bindingCommands, sched)
			err := settingsService.Save(t.Context(), s)

			assert.Error(t, err)
		})
	})
}
