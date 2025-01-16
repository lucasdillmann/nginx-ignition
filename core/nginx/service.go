package nginx

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/nginx/cfgfiles"
	"dillmann.com.br/nginx-ignition/core/settings"
	"github.com/google/uuid"
	"sync"
)

type service struct {
	configFilesManager *cfgfiles.Facade
	processManager     *processManager
	semaphore          *semaphore
	logReader          *logReader
	logRotator         *logRotator
	mu                 sync.Mutex
}

func newService(
	configuration *configuration.Configuration,
	settingsRepository *settings.Repository,
	hostRepository *host.Repository,
	configFilesManager *cfgfiles.Facade,
) (*service, error) {
	pManager, err := newProcessManager(configuration)
	if err != nil {
		return nil, err
	}

	return &service{
		configFilesManager: configFilesManager,
		processManager:     pManager,
		semaphore:          newSemaphore(),
		logReader:          newLogReader(configuration),
		logRotator:         newLogRotator(configuration, settingsRepository, hostRepository, pManager),
	}, nil
}

func (s *service) reload(failIfNotRunning bool) error {
	if failIfNotRunning && s.semaphore.currentState() != runningState {
		return core_error.New("nginx is not running", false)
	}

	return s.semaphore.changeState(runningState, func() error {
		if err := s.configFilesManager.ReplaceConfigurationFiles(); err != nil {
			return err
		}

		return s.processManager.sendReloadSignal()
	})
}

func (s *service) start() error {
	if s.semaphore.currentState() == runningState {
		return nil
	}

	pid, err := s.processManager.currentPid()
	if err != nil {
		return err
	}

	if pid != 0 {
		log.Warnf("nginx seems to be already running with PID %d, trying to reload it instead", pid)
		return s.reload(false)
	}

	return s.semaphore.changeState(runningState, func() error {
		if err := s.configFilesManager.ReplaceConfigurationFiles(); err != nil {
			return err
		}

		return s.processManager.start()
	})
}

func (s *service) stop(_ *int) error {
	if s.semaphore.currentState() == stoppedState {
		return nil
	}

	return s.semaphore.changeState(stoppedState, func() error {
		return s.processManager.sendStopSignal()
	})
}

func (s *service) isRunning() bool {
	return s.semaphore.currentState() == runningState
}

func (s *service) getHostLogs(hostId uuid.UUID, qualifier string, lines int) ([]string, error) {
	return s.logReader.read("host-"+hostId.String()+"."+qualifier+".log", lines)
}

func (s *service) getMainLogs(lines int) ([]string, error) {
	return s.logReader.read("main.log", lines)
}

func (s *service) rotateLogs() error {
	return s.logRotator.rotate()
}
