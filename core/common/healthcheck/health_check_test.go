package healthcheck

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_HealthCheck_Register(t *testing.T) {
	t.Run("registers provider", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		hc := New()
		provider := NewMockProvider(ctrl)
		provider.EXPECT().ID().Return("test-provider")
		provider.EXPECT().Check(ctx).Return(nil)

		hc.Register(provider)

		status := hc.Status(ctx)
		assert.Len(t, status.Details, 1)
	})
}

func Test_HealthCheck_Status(t *testing.T) {
	t.Run("returns healthy when all providers healthy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		hc := New()

		provider1 := NewMockProvider(ctrl)
		provider1.EXPECT().ID().Return("provider1")
		provider1.EXPECT().Check(ctx).Return(nil)

		provider2 := NewMockProvider(ctrl)
		provider2.EXPECT().ID().Return("provider2")
		provider2.EXPECT().Check(ctx).Return(nil)

		hc.Register(provider1)
		hc.Register(provider2)

		status := hc.Status(ctx)

		assert.True(t, status.Healthy)
		assert.Len(t, status.Details, 2)
		assert.NoError(t, status.Details[0].Error)
		assert.NoError(t, status.Details[1].Error)
	})

	t.Run("returns unhealthy when any provider fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		hc := New()

		provider1 := NewMockProvider(ctrl)
		provider1.EXPECT().ID().Return("provider1")
		provider1.EXPECT().Check(ctx).Return(nil)

		provider2 := NewMockProvider(ctrl)
		provider2.EXPECT().ID().Return("provider2")
		provider2.EXPECT().Check(ctx).Return(errors.New("provider error"))

		hc.Register(provider1)
		hc.Register(provider2)

		status := hc.Status(ctx)

		assert.False(t, status.Healthy)
		assert.Len(t, status.Details, 2)
		assert.NoError(t, status.Details[0].Error)
		assert.Error(t, status.Details[1].Error)
	})
}
