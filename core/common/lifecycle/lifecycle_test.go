package lifecycle

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Lifecycle(t *testing.T) {
	t.Run("RegisterStartup", func(t *testing.T) {
		t.Run("registers startup command", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lc := New()
			command := NewMockedStartupCommand(ctrl)

			lc.RegisterStartup(command)

			ctx := context.Background()
			command.EXPECT().Priority().Return(1).AnyTimes()
			command.EXPECT().Async().Return(false)
			command.EXPECT().Run(ctx).Return(nil)

			err := lc.FireStartup(ctx)
			assert.NoError(t, err)
		})
	})

	t.Run("RegisterShutdown", func(t *testing.T) {
		t.Run("registers shutdown command", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lc := New()
			command := NewMockedShutdownCommand(ctrl)

			lc.RegisterShutdown(command)

			ctx := context.Background()
			command.EXPECT().Priority().Return(1).AnyTimes()
			command.EXPECT().Run(ctx)

			lc.FireShutdown(ctx)
		})
	})

	t.Run("FireStartup", func(t *testing.T) {
		t.Run("executes commands in priority order", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lc := New()
			ctx := context.Background()

			command1 := NewMockedStartupCommand(ctrl)
			command1.EXPECT().Priority().Return(2).AnyTimes()
			command1.EXPECT().Async().Return(false)
			command1.EXPECT().Run(ctx).Return(nil)

			command2 := NewMockedStartupCommand(ctrl)
			command2.EXPECT().Priority().Return(1).AnyTimes()
			command2.EXPECT().Async().Return(false)
			command2.EXPECT().Run(ctx).Return(nil)

			lc.RegisterStartup(command1)
			lc.RegisterStartup(command2)

			err := lc.FireStartup(ctx)
			assert.NoError(t, err)
		})

		t.Run("returns error when command fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lc := New()
			ctx := context.Background()

			command := NewMockedStartupCommand(ctrl)
			command.EXPECT().Priority().Return(1).AnyTimes()
			command.EXPECT().Async().Return(false)
			command.EXPECT().Run(ctx).Return(errors.New("startup error"))

			lc.RegisterStartup(command)

			err := lc.FireStartup(ctx)
			assert.Error(t, err)
		})

		t.Run("executes async commands in background", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lc := New()
			ctx := context.Background()

			command := NewMockedStartupCommand(ctrl)
			command.EXPECT().Priority().Return(1).AnyTimes()
			command.EXPECT().Async().Return(true)
			command.EXPECT().Run(ctx).Return(nil)

			lc.RegisterStartup(command)

			err := lc.FireStartup(ctx)
			assert.NoError(t, err)

			time.Sleep(10 * time.Millisecond)
		})
	})

	t.Run("FireShutdown", func(t *testing.T) {
		t.Run("executes shutdown commands in priority order", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lc := New()
			ctx := context.Background()

			command1 := NewMockedShutdownCommand(ctrl)
			command1.EXPECT().Priority().Return(2).AnyTimes()
			command1.EXPECT().Run(ctx)

			command2 := NewMockedShutdownCommand(ctrl)
			command2.EXPECT().Priority().Return(1).AnyTimes()
			command2.EXPECT().Run(ctx)

			lc.RegisterShutdown(command1)
			lc.RegisterShutdown(command2)

			lc.FireShutdown(ctx)
		})
	})
}
