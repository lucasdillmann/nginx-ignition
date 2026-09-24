package nginx

import (
	"archive/zip"
	"bytes"
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/broadcast"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/logline"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
	cfgfiles2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/nginx/cfgfiles"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/vpn"
)

type service struct {
	configFilesManager *cfgfiles2.Facade
	processManager     *processManager
	semaphore          *semaphore
	logReader          *logReader
	logRotator         *logRotator
	vpnManager         *vpnManager
	settingsCommands   settings.Commands
	statsClient        *http.Client
}

func newService(
	cfg *configuration.Configuration,
	hostCommands host.Commands,
	configFilesManager *cfgfiles2.Facade,
	vpnCommands vpn.Commands,
	settingsCommands settings.Commands,
	certificateCommands certificate.Commands,
) (*service, error) {
	pManager, err := newProcessManager(cfg)
	if err != nil {
		return nil, err
	}

	vManager := newVpnManager(vpnCommands, settingsCommands, certificateCommands)

	return &service{
		configFilesManager: configFilesManager,
		processManager:     pManager,
		vpnManager:         vManager,
		settingsCommands:   settingsCommands,
		semaphore:          newSemaphore(),
		logReader:          newLogReader(cfg),
		logRotator:         newLogRotator(cfg, settingsCommands, hostCommands, pManager),
		statsClient:        buildStatsClient(pManager.configPath),
	}, nil
}

func (s *service) Reload(ctx context.Context, failIfNotRunning bool) error {
	if failIfNotRunning && s.semaphore.currentState() != runningState {
		return coreerror.New(i18n.M(ctx, i18n.K.CoreNginxNotRunning), false)
	}

	supportedFeatures, err := s.resolveSupportedFeatures(ctx)
	if err != nil {
		return err
	}

	return s.semaphore.changeState(runningState, func() error {
		hosts, _, err := s.configFilesManager.ReplaceConfigurationFiles(ctx, supportedFeatures)
		if err != nil {
			return err
		}

		err = s.processManager.sendReloadSignal()
		if err != nil {
			return err
		}

		return s.vpnManager.reload(ctx, hosts)
	})
}

func (s *service) Start(ctx context.Context) error {
	if s.semaphore.currentState() == runningState {
		return nil
	}

	pid, err := s.processManager.currentPid()
	if err != nil {
		return err
	}

	if pid != 0 {
		log.Warnf("nginx seems to be already running with PID %d, trying to reload it instead", pid)
		return s.Reload(ctx, false)
	}

	supportedFeatures, err := s.resolveSupportedFeatures(ctx)
	if err != nil {
		return err
	}

	return s.semaphore.changeState(runningState, func() error {
		hosts, _, err := s.configFilesManager.ReplaceConfigurationFiles(ctx, supportedFeatures)
		if err != nil {
			return err
		}

		err = s.processManager.start()
		if err != nil {
			return err
		}

		return s.vpnManager.start(ctx, hosts)
	})
}

func (s *service) Stop(ctx context.Context) error {
	if s.semaphore.currentState() == stoppedState {
		return nil
	}

	return s.semaphore.changeState(stoppedState, func() error {
		if err := s.vpnManager.stop(ctx); err != nil {
			return err
		}

		return s.processManager.sendStopSignal()
	})
}

func (s *service) GetStatus(_ context.Context) Status {
	running := s.semaphore.currentState() == runningState
	status := Status{Running: running}
	if !running {
		return status
	}

	uptime, err := s.processManager.uptimeSeconds()
	if err != nil {
		log.Warnf("unable to resolve nginx uptime: %v", err)
		return status
	}

	status.UptimeSeconds = &uptime
	return status
}

func (s *service) GetHostLogs(
	ctx context.Context,
	hostID uuid.UUID,
	qualifier string,
	lines int,
	search *LogSearch,
) ([]logline.LogLine, error) {
	return s.readLogs(ctx, lines, "host-"+hostID.String()+"."+qualifier+".log", search)
}

func (s *service) GetMainLogs(
	ctx context.Context,
	lines int,
	search *LogSearch,
) ([]logline.LogLine, error) {
	return s.readLogs(ctx, lines, "main.log", search)
}

func (s *service) readLogs(
	ctx context.Context,
	lines int,
	fileName string,
	search *LogSearch,
) ([]logline.LogLine, error) {
	output, err := s.logReader.read(ctx, fileName)
	if err != nil {
		return nil, err
	}

	if search != nil {
		output, err = logline.Search(output, search.Query, search.SurroundingLines)
		if err != nil {
			return nil, err
		}
	}

	if len(output) > lines {
		output = output[len(output)-lines:]
	}

	return output, nil
}

func (s *service) rotateLogs(ctx context.Context) error {
	return s.logRotator.rotate(ctx)
}

func (s *service) attachListeners() {
	channel := broadcast.Listen("core:nginx:reload")
	for range channel {
		err := s.Reload(<-channel, false)
		if err != nil {
			log.Warnf("Failed to reload nginx: %v", err)
		}
	}
}

func (s *service) GetConfigFiles(
	ctx context.Context,
	input GetConfigFilesInput,
) ([]byte, error) {
	paths := &cfgfiles2.Paths{
		Base:   input.BasePath,
		Config: input.ConfigPath,
		Logs:   input.LogPath,
		Cache:  input.CachePath,
		Temp:   input.TempPath,
	}

	supportedFeatures, err := s.resolveSupportedFeatures(ctx)
	if err != nil {
		return nil, err
	}

	configFiles, _, _, err := s.configFilesManager.GetConfigurationFiles(
		ctx,
		paths,
		supportedFeatures,
	)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)

	//nolint:errcheck
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

func (s *service) resolveSupportedFeatures(
	ctx context.Context,
) (*cfgfiles2.SupportedFeatures, error) {
	metadata, err := s.GetMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return &cfgfiles2.SupportedFeatures{
		TLSSNI:      cfgfiles2.SupportType(metadata.SNISupportType()),
		RunCodeType: cfgfiles2.SupportType(metadata.RunCodeSupportType()),
		StreamType:  cfgfiles2.SupportType(metadata.StreamSupportType()),
		StatsType:   cfgfiles2.SupportType(metadata.StatsSupportType()),
		GRPCType:    cfgfiles2.SupportType(metadata.GRPCSupportType()),
	}, nil
}
