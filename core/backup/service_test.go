package backup

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Service(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("returns backup when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expected := &Backup{
				FileName:    "backup.db",
				ContentType: "application/octet-stream",
				Contents:    []byte("test content"),
			}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Get(ctx).Return(expected, nil)

			svc := newCommands(repo)
			result, err := svc.Get(ctx)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expectedErr := errors.New("repository error")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Get(ctx).Return(nil, expectedErr)

			svc := newCommands(repo)
			result, err := svc.Get(ctx)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, expectedErr, err)
		})
	})
}
