package nginx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	settings2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/settings"
)

func Test_logRotationTask(t *testing.T) {
	t.Run("Schedule", func(t *testing.T) {
		t.Run("converts minutes to duration correctly", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := settings2.NewMockedCommands(ctrl)
			repo.EXPECT().Get(t.Context()).Return(&settings2.Settings{
				LogRotation: &settings2.LogRotationSettings{
					Enabled:           true,
					IntervalUnit:      settings2.MinutesTimeUnit,
					IntervalUnitCount: 30,
				},
			}, nil)

			task := &logRotationTask{
				settingsCommands: repo,
			}
			schedule, err := task.Schedule(t.Context())

			assert.NoError(t, err)
			assert.True(t, schedule.Enabled)
			assert.Equal(t, 30*time.Minute, schedule.Interval)
		})

		t.Run("converts hours to duration correctly", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := settings2.NewMockedCommands(ctrl)
			repo.EXPECT().Get(t.Context()).Return(&settings2.Settings{
				LogRotation: &settings2.LogRotationSettings{
					Enabled:           true,
					IntervalUnit:      settings2.HoursTimeUnit,
					IntervalUnitCount: 2,
				},
			}, nil)

			task := &logRotationTask{
				settingsCommands: repo,
			}
			schedule, err := task.Schedule(t.Context())

			assert.NoError(t, err)
			assert.Equal(t, 2*time.Hour, schedule.Interval)
		})

		t.Run("converts days to duration correctly", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := settings2.NewMockedCommands(ctrl)
			repo.EXPECT().Get(t.Context()).Return(&settings2.Settings{
				LogRotation: &settings2.LogRotationSettings{
					Enabled:           false,
					IntervalUnit:      settings2.DaysTimeUnit,
					IntervalUnitCount: 1,
				},
			}, nil)

			task := &logRotationTask{
				settingsCommands: repo,
			}
			schedule, err := task.Schedule(t.Context())

			assert.NoError(t, err)
			assert.False(t, schedule.Enabled)
			assert.Equal(t, 24*time.Hour, schedule.Interval)
		})

		t.Run("returns error for invalid unit", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := settings2.NewMockedCommands(ctrl)
			repo.EXPECT().Get(t.Context()).Return(&settings2.Settings{
				LogRotation: &settings2.LogRotationSettings{
					IntervalUnit: "invalid",
				},
			}, nil)

			task := &logRotationTask{
				settingsCommands: repo,
			}
			_, err := task.Schedule(t.Context())

			assert.Error(t, err)
		})

		t.Run("returns error when settings retrieval fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := settings2.NewMockedCommands(ctrl)
			repo.EXPECT().Get(t.Context()).Return(nil, assert.AnError)

			task := &logRotationTask{
				settingsCommands: repo,
			}
			_, err := task.Schedule(t.Context())

			assert.Error(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})
}
