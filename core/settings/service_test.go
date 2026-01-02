package settings

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
)

func buildScheduler() *scheduler.Scheduler {
	s := &scheduler.Scheduler{}
	return s
}

func Test_Service_Get(t *testing.T) {
	t.Run("returns settings when found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		expected := &Settings{}

		repo := NewMockRepository(ctrl)
		repo.EXPECT().Get(ctx).Return(expected, nil)

		bindingCommands := &binding.Commands{}
		sched := buildScheduler()

		svc := newService(repo, bindingCommands, sched)
		result, err := svc.get(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		expectedErr := errors.New("repository error")

		repo := NewMockRepository(ctrl)
		repo.EXPECT().Get(ctx).Return(nil, expectedErr)

		bindingCommands := &binding.Commands{}
		sched := buildScheduler()

		svc := newService(repo, bindingCommands, sched)
		result, err := svc.get(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func Test_Service_Save(t *testing.T) {
	t.Run("invalid settings returns validation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		settings := &Settings{
			Nginx: &NginxSettings{
				DefaultContentType: "",
				RuntimeUser:        "nginx",
				Timeouts: &NginxTimeoutsSettings{
					Read:      60,
					Send:      60,
					Connect:   60,
					Keepalive: 65,
				},
				WorkerProcesses:   1,
				WorkerConnections: 1024,
				MaximumBodySizeMb: 1,
			},
			LogRotation: &LogRotationSettings{
				IntervalUnitCount: 1,
				MaximumLines:      1000,
			},
			CertificateAutoRenew: &CertificateAutoRenewSettings{
				IntervalUnitCount: 1,
			},
		}

		repo := NewMockRepository(ctrl)
		bindingCommands := &binding.Commands{}
		sched := buildScheduler()

		svc := newService(repo, bindingCommands, sched)
		err := svc.save(ctx, settings)

		assert.Error(t, err)
	})
}
