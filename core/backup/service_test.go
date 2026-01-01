package backup

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_Get(t *testing.T) {
	t.Run("returns backup when found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		expected := &Backup{
			FileName:    "backup.db",
			ContentType: "application/octet-stream",
			Contents:    []byte("test content"),
		}

		repo := NewMockRepository(ctrl)
		repo.EXPECT().Get(ctx).Return(expected, nil)

		svc := newService(repo)
		result, err := svc.get(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		expectedErr := errors.New("repository error")

		repo := NewMockRepository(ctrl)
		repo.EXPECT().Get(ctx).Return(nil, expectedErr)

		svc := newService(repo)
		result, err := svc.get(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}
