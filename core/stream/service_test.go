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

func TestService_Save(t *testing.T) {
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

		repo := NewMockRepository(ctrl)
		repo.EXPECT().Save(ctx, stream).Return(nil)

		svc := newService(repo)
		err := svc.save(ctx, stream)

		assert.NoError(t, err)
	})

	t.Run("invalid stream returns validation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		stream := &Stream{
			Name: "",
		}

		repo := NewMockRepository(ctrl)
		svc := newService(repo)
		err := svc.save(ctx, stream)

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
		repo := NewMockRepository(ctrl)
		repo.EXPECT().Save(ctx, stream).Return(expectedErr)

		svc := newService(repo)
		err := svc.save(ctx, stream)

		assert.Equal(t, expectedErr, err)
	})
}

func TestService_DeleteByID(t *testing.T) {
	t.Run("deletes successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()

		repo := NewMockRepository(ctrl)
		repo.EXPECT().DeleteByID(ctx, id).Return(nil)

		svc := newService(repo)
		err := svc.deleteByID(ctx, id)

		assert.NoError(t, err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		expectedErr := errors.New("delete failed")

		repo := NewMockRepository(ctrl)
		repo.EXPECT().DeleteByID(ctx, id).Return(expectedErr)

		svc := newService(repo)
		err := svc.deleteByID(ctx, id)

		assert.Equal(t, expectedErr, err)
	})
}

func TestService_GetByID(t *testing.T) {
	t.Run("returns stream when found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		expected := &Stream{
			ID:   id,
			Name: "test",
		}

		repo := NewMockRepository(ctrl)
		repo.EXPECT().FindByID(ctx, id).Return(expected, nil)

		svc := newService(repo)
		result, err := svc.getByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestService_List(t *testing.T) {
	t.Run("returns paginated results", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		expectedPage := pagination.New(1, 10, 1, []Stream{
			{Name: "test"},
		})
		searchTerms := "test"

		repo := NewMockRepository(ctrl)
		repo.EXPECT().FindPage(ctx, 10, 1, &searchTerms).Return(expectedPage, nil)

		svc := newService(repo)
		result, err := svc.list(ctx, 10, 1, &searchTerms)

		assert.NoError(t, err)
		assert.Equal(t, expectedPage, result)
	})
}

func TestService_ExistsByID(t *testing.T) {
	t.Run("returns true when exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()

		repo := NewMockRepository(ctrl)
		repo.EXPECT().ExistsByID(ctx, id).Return(true, nil)

		svc := newService(repo)
		exists, err := svc.existsByID(ctx, id)

		assert.NoError(t, err)
		assert.True(t, exists)
	})
}
