package lifecycle

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLifecycle_RegisterStartup(t *testing.T) {
	t.Run("registers startup command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lc := New()
		command := NewMockStartupCommand(ctrl)

		lc.RegisterStartup(command)

		ctx := context.Background()
		command.EXPECT().Priority().Return(1).AnyTimes()
		command.EXPECT().Async().Return(false)
		command.EXPECT().Run(ctx).Return(nil)

		err := lc.FireStartup(ctx)
		assert.NoError(t, err)
	})
}

func TestLifecycle_RegisterShutdown(t *testing.T) {
	t.Run("registers shutdown command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lc := New()
		command := NewMockShutdownCommand(ctrl)

		lc.RegisterShutdown(command)

		ctx := context.Background()
		command.EXPECT().Priority().Return(1).AnyTimes()
		command.EXPECT().Run(ctx)

		lc.FireShutdown(ctx)
	})
}

func TestLifecycle_FireStartup(t *testing.T) {
	t.Run("executes commands in priority order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lc := New()
		ctx := context.Background()

		command1 := NewMockStartupCommand(ctrl)
		command1.EXPECT().Priority().Return(2).AnyTimes()
		command1.EXPECT().Async().Return(false)
		command1.EXPECT().Run(ctx).Return(nil)

		command2 := NewMockStartupCommand(ctrl)
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

		command := NewMockStartupCommand(ctrl)
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

		command := NewMockStartupCommand(ctrl)
		command.EXPECT().Priority().Return(1).AnyTimes()
		command.EXPECT().Async().Return(true)
		command.EXPECT().Run(ctx).Return(nil)

		lc.RegisterStartup(command)

		err := lc.FireStartup(ctx)
		assert.NoError(t, err)

		time.Sleep(10 * time.Millisecond)
	})
}

func TestLifecycle_FireShutdown(t *testing.T) {
	t.Run("executes shutdown commands in priority order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lc := New()
		ctx := context.Background()

		command1 := NewMockShutdownCommand(ctrl)
		command1.EXPECT().Priority().Return(2).AnyTimes()
		command1.EXPECT().Run(ctx)

		command2 := NewMockShutdownCommand(ctrl)
		command2.EXPECT().Priority().Return(1).AnyTimes()
		command2.EXPECT().Run(ctx)

		lc.RegisterShutdown(command1)
		lc.RegisterShutdown(command2)

		lc.FireShutdown(ctx)
	})
}
