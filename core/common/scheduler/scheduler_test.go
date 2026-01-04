package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_BuildScheduler(t *testing.T) {
	t.Run("builds scheduler", func(t *testing.T) {
		sched := buildScheduler()

		assert.NotNil(t, sched)
		assert.NotNil(t, sched.tickers)
		assert.False(t, sched.stopped)
		assert.False(t, sched.started)
	})
}

func Test_Scheduler(t *testing.T) {
	t.Run("Register", func(t *testing.T) {
		t.Run("registers task when not started", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sched := buildScheduler()
			task := NewMockedTask(ctrl)

			err := sched.Register(context.Background(), task)

			assert.NoError(t, err)
			assert.Contains(t, sched.tickers, task)
		})

		t.Run("starts task when scheduler already started", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			sched := buildScheduler()
			sched.started = true

			task := NewMockedTask(ctrl)
			task.EXPECT().Schedule(ctx).Return(&Schedule{
				Enabled:  true,
				Interval: time.Second,
			}, nil)
			task.EXPECT().OnScheduleStarted(ctx)

			err := sched.Register(ctx, task)

			assert.NoError(t, err)
		})

		t.Run("returns error when stopped", func(t *testing.T) {
			sched := buildScheduler()
			sched.stopped = true
			task := NewMockedTask(gomock.NewController(t))

			err := sched.Register(context.Background(), task)

			assert.Error(t, err)
		})
	})

	t.Run("start", func(t *testing.T) {
		t.Run("starts all registered tasks", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			sched := buildScheduler()

			task := NewMockedTask(ctrl)
			task.EXPECT().Schedule(ctx).Return(&Schedule{
				Enabled:  true,
				Interval: time.Second,
			}, nil)
			task.EXPECT().OnScheduleStarted(ctx)

			sched.tickers[task] = time.NewTicker(time.Hour)

			err := sched.start(ctx)

			assert.NoError(t, err)
			assert.True(t, sched.started)
		})

		t.Run("returns error when already started", func(t *testing.T) {
			sched := buildScheduler()
			sched.started = true

			err := sched.start(context.Background())

			assert.Error(t, err)
		})

		t.Run("returns error when stopped", func(t *testing.T) {
			sched := buildScheduler()
			sched.stopped = true

			err := sched.start(context.Background())

			assert.Error(t, err)
		})

		t.Run("returns error when task schedule fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			sched := buildScheduler()

			task := NewMockedTask(ctrl)
			task.EXPECT().Schedule(ctx).Return(nil, errors.New("schedule error"))

			sched.tickers[task] = time.NewTicker(time.Hour)

			err := sched.start(ctx)

			assert.Error(t, err)
		})
	})

	t.Run("stop", func(t *testing.T) {
		t.Run("stops all tasks", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sched := buildScheduler()
			task := NewMockedTask(ctrl)
			sched.tickers[task] = time.NewTicker(time.Second)

			sched.stop()

			assert.True(t, sched.stopped)
			assert.Empty(t, sched.tickers)
		})
	})

	t.Run("Reload", func(t *testing.T) {
		t.Run("reloads all tasks", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			sched := buildScheduler()

			task := NewMockedTask(ctrl)
			task.EXPECT().Schedule(ctx).Return(&Schedule{
				Enabled:  true,
				Interval: time.Minute,
			}, nil)
			task.EXPECT().OnScheduleStarted(ctx)

			sched.tickers[task] = time.NewTicker(time.Hour)

			err := sched.Reload(ctx)

			assert.NoError(t, err)
		})

		t.Run("returns error when stopped", func(t *testing.T) {
			sched := buildScheduler()
			sched.stopped = true

			err := sched.Reload(context.Background())

			assert.Error(t, err)
		})

		t.Run("returns error when task schedule fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			sched := buildScheduler()

			task := NewMockedTask(ctrl)
			task.EXPECT().Schedule(ctx).Return(nil, errors.New("schedule error"))

			sched.tickers[task] = time.NewTicker(time.Hour)

			err := sched.Reload(ctx)

			assert.Error(t, err)
		})
	})
}
