package nginx

import (
	"archive/zip"
	"bytes"
	"context"
	"sync"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/broadcast"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/nginx/cfgfiles"
	"dillmann.com.br/nginx-ignition/core/settings"
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
	settingsRepository settings.Repository,
	hostRepository host.Repository,
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

func (s *service) reload(ctx context.Context, failIfNotRunning bool) error {
	if failIfNotRunning && s.semaphore.currentState() != runningState {
		return core_error.New("nginx is not running", false)
	}

	supportedFeatures, err := s.resolveSupportedFeatures(ctx)
	if err != nil {
		return err
	}

	return s.semaphore.changeState(runningState, func() error {
		if err := s.configFilesManager.ReplaceConfigurationFiles(ctx, supportedFeatures); err != nil {
			return err
		}

		return s.processManager.sendReloadSignal()
	})
}

func (s *service) start(ctx context.Context) error {
	if s.semaphore.currentState() == runningState {
		return nil
	}

	pid, err := s.processManager.currentPid()
	if err != nil {
		return err
	}

	if pid != 0 {
		log.Warnf("nginx seems to be already running with PID %d, trying to reload it instead", pid)
		return s.reload(ctx, false)
	}

	supportedFeatures, err := s.resolveSupportedFeatures(ctx)
	if err != nil {
		return err
	}

	return s.semaphore.changeState(runningState, func() error {
		if err := s.configFilesManager.ReplaceConfigurationFiles(ctx, supportedFeatures); err != nil {
			return err
		}

		return s.processManager.start()
	})
}

func (s *service) stop(_ context.Context) error {
	if s.semaphore.currentState() == stoppedState {
		return nil
	}

	return s.semaphore.changeState(stoppedState, func() error {
		return s.processManager.sendStopSignal()
	})
}

func (s *service) isRunning(_ context.Context) bool {
	return s.semaphore.currentState() == runningState
}

func (s *service) getHostLogs(ctx context.Context, hostId uuid.UUID, qualifier string, lines int) ([]string, error) {
	return s.logReader.read(ctx, "host-"+hostId.String()+"."+qualifier+".log", lines)
}

func (s *service) getMainLogs(ctx context.Context, lines int) ([]string, error) {
	return s.logReader.read(ctx, "main.log", lines)
}

func (s *service) rotateLogs(ctx context.Context) error {
	return s.logRotator.rotate(ctx)
}

func (s *service) attachListeners() {
	channel := broadcast.Listen("core:nginx:reload")
	for range channel {
		err := s.reload(<-channel, false)
		if err != nil {
			log.Warnf("Failed to reload nginx: %v", err)
		}
	}
}

func (s *service) getConfigFilesZipFile(
	ctx context.Context,
	basePath, configPath, logPath string,
) ([]byte, error) {
	paths := &cfgfiles.Paths{
		Base:   basePath,
		Config: configPath,
		Logs:   logPath,
	}

	supportedFeatures, err := s.resolveSupportedFeatures(ctx)
	if err != nil {
		return nil, err
	}

	configFiles, _, _, err := s.configFilesManager.GetConfigurationFiles(ctx, paths, supportedFeatures)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)
	defer zipWriter.Close()

	for _, file := range configFiles {
		itemWriter, err := zipWriter.Create(file.Name)
		if err != nil {
			return nil, err
		}

		if _, err := itemWriter.Write([]byte(file.FormattedContents())); err != nil {
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *service) resolveSupportedFeatures(ctx context.Context) (*cfgfiles.SupportedFeatures, error) {
	metadata, err := s.getMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return &cfgfiles.SupportedFeatures{
		TLSSNI:      cfgfiles.SupportType(metadata.SNISupportType()),
		RunCodeType: cfgfiles.SupportType(metadata.RunCodeSupportType()),
		StreamType:  cfgfiles.SupportType(metadata.StreamSupportType()),
	}, nil
}
