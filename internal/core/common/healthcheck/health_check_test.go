package healthcheck

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_HealthCheck(t *testing.T) {
	t.Run("Register", func(t *testing.T) {
		t.Run("registers provider", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			hc := New()
			provider := NewMockedProvider(ctrl)
			provider.EXPECT().ID().Return("test-provider")
			provider.EXPECT().Check(t.Context()).Return(nil)

			hc.Register(provider)

			status := hc.Status(t.Context())
			assert.Len(t, status.Details, 1)
		})
	})

	t.Run("Status", func(t *testing.T) {
		t.Run("returns healthy when all providers healthy", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			hc := New()

			provider1 := NewMockedProvider(ctrl)
			provider1.EXPECT().ID().Return("provider1")
			provider1.EXPECT().Check(t.Context()).Return(nil)

			provider2 := NewMockedProvider(ctrl)
			provider2.EXPECT().ID().Return("provider2")
			provider2.EXPECT().Check(t.Context()).Return(nil)

			hc.Register(provider1)
			hc.Register(provider2)

			status := hc.Status(t.Context())

			assert.True(t, status.Healthy)
			assert.Len(t, status.Details, 2)
			assert.NoError(t, status.Details[0].Error)
			assert.NoError(t, status.Details[1].Error)
		})

		t.Run("returns unhealthy when any provider fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			hc := New()

			provider1 := NewMockedProvider(ctrl)
			provider1.EXPECT().ID().Return("provider1")
			provider1.EXPECT().Check(t.Context()).Return(nil)

			provider2 := NewMockedProvider(ctrl)
			provider2.EXPECT().ID().Return("provider2")
			provider2.EXPECT().Check(t.Context()).Return(errors.New("provider error"))

			hc.Register(provider1)
			hc.Register(provider2)

			status := hc.Status(t.Context())

			assert.False(t, status.Healthy)
			assert.Len(t, status.Details, 2)
			assert.NoError(t, status.Details[0].Error)
			assert.Error(t, status.Details[1].Error)
		})
	})
}
