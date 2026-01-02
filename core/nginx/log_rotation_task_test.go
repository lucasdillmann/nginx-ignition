package nginx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_LogRotationTask_Schedule(t *testing.T) {
	ctx := context.Background()

	t.Run("converts minutes to duration correctly", func(t *testing.T) {
		repo := &settings.Commands{
			Get: func(_ context.Context) (*settings.Settings, error) {
				return &settings.Settings{
					LogRotation: &settings.LogRotationSettings{
						Enabled:           true,
						IntervalUnit:      settings.MinutesTimeUnit,
						IntervalUnitCount: 30,
					},
				}, nil
			},
		}

		task := &logRotationTask{
			settingsCommands: repo,
		}
		schedule, err := task.Schedule(ctx)

		assert.NoError(t, err)
		assert.True(t, schedule.Enabled)
		assert.Equal(t, 30*time.Minute, schedule.Interval)
	})

	t.Run("converts hours to duration correctly", func(t *testing.T) {
		repo := &settings.Commands{
			Get: func(_ context.Context) (*settings.Settings, error) {
				return &settings.Settings{
					LogRotation: &settings.LogRotationSettings{
						Enabled:           true,
						IntervalUnit:      settings.HoursTimeUnit,
						IntervalUnitCount: 2,
					},
				}, nil
			},
		}

		task := &logRotationTask{
			settingsCommands: repo,
		}
		schedule, err := task.Schedule(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 2*time.Hour, schedule.Interval)
	})

	t.Run("converts days to duration correctly", func(t *testing.T) {
		repo := &settings.Commands{
			Get: func(_ context.Context) (*settings.Settings, error) {
				return &settings.Settings{
					LogRotation: &settings.LogRotationSettings{
						Enabled:           false,
						IntervalUnit:      settings.DaysTimeUnit,
						IntervalUnitCount: 1,
					},
				}, nil
			},
		}

		task := &logRotationTask{
			settingsCommands: repo,
		}
		schedule, err := task.Schedule(ctx)

		assert.NoError(t, err)
		assert.False(t, schedule.Enabled)
		assert.Equal(t, 24*time.Hour, schedule.Interval)
	})

	t.Run("returns error for invalid unit", func(t *testing.T) {
		repo := &settings.Commands{
			Get: func(_ context.Context) (*settings.Settings, error) {
				return &settings.Settings{
					LogRotation: &settings.LogRotationSettings{
						IntervalUnit: "invalid",
					},
				}, nil
			},
		}

		task := &logRotationTask{
			settingsCommands: repo,
		}
		_, err := task.Schedule(ctx)

		assert.Error(t, err)
	})

	t.Run("returns error when settings retrieval fails", func(t *testing.T) {
		repo := &settings.Commands{
			Get: func(_ context.Context) (*settings.Settings, error) {
				return nil, assert.AnError
			},
		}

		task := &logRotationTask{
			settingsCommands: repo,
		}
		_, err := task.Schedule(ctx)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}
