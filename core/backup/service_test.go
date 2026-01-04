package backup

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_service(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("returns backup when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expected := newBackup()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().Get(ctx).Return(expected, nil)

			backupService := newCommands(repository)
			result, err := backupService.Get(ctx)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expectedErr := errors.New("repository error")

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().Get(ctx).Return(nil, expectedErr)

			backupService := newCommands(repository)
			result, err := backupService.Get(ctx)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, expectedErr, err)
		})
	})
}
