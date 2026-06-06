package notification

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/scheduler"
)

func Test_deliveryTask(t *testing.T) {
	t.Run("Schedule", func(t *testing.T) {
		t.Run("returns enabled schedule with 30 second interval", func(t *testing.T) {
			task := &deliveryTask{}

			schedule, err := task.Schedule(t.Context())

			assert.NoError(t, err)
			assert.True(t, schedule.Enabled)
			assert.Equal(t, 30*time.Second, schedule.Interval)
		})
	})

	t.Run("Run", func(t *testing.T) {
		t.Run("calls ProcessPendingDeliveries", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := NewMockedCommands(ctrl)
			commands.EXPECT().ProcessPendingDeliveries(t.Context()).Return(nil)

			task := deliveryTask{commands: commands}

			err := task.Run(t.Context())

			assert.NoError(t, err)
		})

		t.Run("propagates ProcessPendingDeliveries errors", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := NewMockedCommands(ctrl)
			commands.EXPECT().ProcessPendingDeliveries(t.Context()).Return(assert.AnError)

			task := deliveryTask{commands: commands}

			err := task.Run(t.Context())

			assert.Error(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("registerScheduledTask", func(t *testing.T) {
		t.Run("wires delivery task to commands", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := NewMockedCommands(ctrl)
			commands.EXPECT().ProcessPendingDeliveries(t.Context()).Return(nil)

			task := deliveryTask{commands: commands}
			err := task.Run(t.Context())

			assert.NoError(t, err)
		})

		t.Run("registers task with scheduler", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := NewMockedCommands(ctrl)
			sched := scheduler.New()

			err := registerScheduledTask(t.Context(), commands, sched)

			assert.NoError(t, err)
		})
	})
}
