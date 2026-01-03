package stream

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func Test_Service(t *testing.T) {
	t.Run("Save", func(t *testing.T) {
		t.Run("valid stream saves successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			port := 8080
			stream := &Stream{
				Name: "test",
				Type: SimpleType,
				Binding: Address{
					Protocol: TCPProtocol,
					Address:  "127.0.0.1",
					Port:     &port,
				},
				DefaultBackend: Backend{
					Address: Address{
						Protocol: TCPProtocol,
						Address:  "127.0.0.1",
						Port:     &port,
					},
				},
			}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Save(ctx, stream).Return(nil)

			svc := newCommands(repo)
			err := svc.Save(ctx, stream)

			assert.NoError(t, err)
		})

		t.Run("invalid stream returns validation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			stream := &Stream{
				Name: "",
			}

			repo := NewMockedRepository(ctrl)
			svc := newCommands(repo)
			err := svc.Save(ctx, stream)

			assert.Error(t, err)
		})

		t.Run("repository error is returned", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			port := 8080
			stream := &Stream{
				Name: "test",
				Type: SimpleType,
				Binding: Address{
					Protocol: TCPProtocol,
					Address:  "127.0.0.1",
					Port:     &port,
				},
				DefaultBackend: Backend{
					Address: Address{
						Protocol: TCPProtocol,
						Address:  "127.0.0.1",
						Port:     &port,
					},
				},
			}

			expectedErr := errors.New("repository error")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Save(ctx, stream).Return(expectedErr)

			svc := newCommands(repo)
			err := svc.Save(ctx, stream)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("deletes successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().DeleteByID(ctx, id).Return(nil)

			svc := newCommands(repo)
			err := svc.Delete(ctx, id)

			assert.NoError(t, err)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expectedErr := errors.New("delete failed")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().DeleteByID(ctx, id).Return(expectedErr)

			svc := newCommands(repo)
			err := svc.Delete(ctx, id)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("returns stream when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expected := &Stream{
				ID:   id,
				Name: "test",
			}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(ctx, id).Return(expected, nil)

			svc := newCommands(repo)
			result, err := svc.Get(ctx, id)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("returns paginated results", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expectedPage := pagination.New(1, 10, 1, []Stream{
				{
					Name: "test",
				},
			})
			searchTerms := "test"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindPage(ctx, 10, 1, &searchTerms).Return(expectedPage, nil)

			svc := newCommands(repo)
			result, err := svc.List(ctx, 10, 1, &searchTerms)

			assert.NoError(t, err)
			assert.Equal(t, expectedPage, result)
		})
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("returns true when exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().ExistsByID(ctx, id).Return(true, nil)

			svc := newCommands(repo)
			exists, err := svc.Exists(ctx, id)

			assert.NoError(t, err)
			assert.True(t, exists)
		})
	})
}
