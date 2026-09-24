package cfgfiles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	cache2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/cache"
	host2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
	settings2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/settings"
	stream2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/stream"
)

func Test_Facade(t *testing.T) {
	t.Run("GetConfigurationFiles", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		paths := &Paths{
			Base: "/etc/nginx/",
		}
		features := &SupportedFeatures{
			StreamType: StaticSupportType,
		}

		t.Run("successfully collects files from providers", func(t *testing.T) {
			hostCmds := host2.NewMockedCommands(ctrl)
			hostCmds.EXPECT().
				GetAllEnabled(t.Context()).
				Return([]host2.Host{
					{
						ID:          uuid.New(),
						DomainNames: []string{"example.com"},
					},
				}, nil)
			streamCmds := stream2.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(t.Context()).Return([]stream2.Stream{}, nil)
			cacheCmds := cache2.NewMockedCommands(ctrl)
			cacheCmds.EXPECT().GetAllInUse(t.Context()).Return([]cache2.Cache{}, nil)

			provider := NewMockedfileProvider(ctrl)
			provider.EXPECT().
				provide(gomock.Any()).
				Return([]File{
					{
						Name:     "test.conf",
						Contents: "test",
					},
				}, nil)

			settingsCmds := settings2.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(t.Context()).Return(&settings2.Settings{}, nil)

			facade := &Facade{
				hostCommands:     hostCmds,
				streamCommands:   streamCmds,
				cacheCommands:    cacheCmds,
				settingsCommands: settingsCmds,
				providers:        []fileProvider{provider},
			}

			configFiles, hosts, streams, err := facade.GetConfigurationFiles(
				t.Context(),
				paths,
				features,
			)

			assert.NoError(t, err)
			assert.Len(t, configFiles, 1)
			assert.Equal(t, "test.conf", configFiles[0].Name)
			assert.Len(t, hosts, 1)
			assert.Len(t, streams, 0)
		})

		t.Run("returns error when hostCommands fails", func(t *testing.T) {
			hostCmds := host2.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(t.Context()).Return(nil, assert.AnError)
			facade := &Facade{hostCommands: hostCmds}
			_, _, _, err := facade.GetConfigurationFiles(t.Context(), paths, features)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("returns error when settingsCommands fails", func(t *testing.T) {
			hostCmds := host2.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(gomock.Any()).Return([]host2.Host{}, nil)
			streamCmds := stream2.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(gomock.Any()).Return([]stream2.Stream{}, nil)
			cacheCmds := cache2.NewMockedCommands(ctrl)
			cacheCmds.EXPECT().GetAllInUse(gomock.Any()).Return([]cache2.Cache{}, nil)
			settingsCmds := settings2.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).Return(nil, assert.AnError)

			facade := &Facade{
				hostCommands:     hostCmds,
				streamCommands:   streamCmds,
				cacheCommands:    cacheCmds,
				settingsCommands: settingsCmds,
			}
			_, _, _, err := facade.GetConfigurationFiles(t.Context(), paths, features)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("returns error when streamCommands fails", func(t *testing.T) {
			hostCmds := host2.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(t.Context()).Return([]host2.Host{}, nil)
			streamCmds := stream2.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(t.Context()).Return(nil, assert.AnError)
			settingsCmds := settings2.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings2.Settings{}, nil).AnyTimes()
			facade := &Facade{
				hostCommands:     hostCmds,
				streamCommands:   streamCmds,
				settingsCommands: settingsCmds,
			}
			_, _, _, err := facade.GetConfigurationFiles(t.Context(), paths, features)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("returns error when cacheCommands fails", func(t *testing.T) {
			hostCmds := host2.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(t.Context()).Return([]host2.Host{}, nil)
			streamCmds := stream2.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(t.Context()).Return([]stream2.Stream{}, nil)
			cacheCmds := cache2.NewMockedCommands(ctrl)
			cacheCmds.EXPECT().GetAllInUse(t.Context()).Return(nil, assert.AnError)
			settingsCmds := settings2.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings2.Settings{}, nil).AnyTimes()
			facade := &Facade{
				hostCommands:     hostCmds,
				streamCommands:   streamCmds,
				cacheCommands:    cacheCmds,
				settingsCommands: settingsCmds,
			}
			_, _, _, err := facade.GetConfigurationFiles(t.Context(), paths, features)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("returns error when provider fails", func(t *testing.T) {
			hostCmds := host2.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(t.Context()).Return([]host2.Host{}, nil)
			streamCmds := stream2.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(t.Context()).Return([]stream2.Stream{}, nil)
			cacheCmds := cache2.NewMockedCommands(ctrl)
			cacheCmds.EXPECT().GetAllInUse(t.Context()).Return([]cache2.Cache{}, nil)

			provider := NewMockedfileProvider(ctrl)
			provider.EXPECT().provide(gomock.Any()).Return(nil, assert.AnError)

			settingsCmds := settings2.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings2.Settings{}, nil).AnyTimes()

			facade := &Facade{
				hostCommands:     hostCmds,
				streamCommands:   streamCmds,
				cacheCommands:    cacheCmds,
				settingsCommands: settingsCmds,
				providers:        []fileProvider{provider},
			}

			_, _, _, err := facade.GetConfigurationFiles(t.Context(), paths, features)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("collects files from multiple providers", func(t *testing.T) {
			hostCmds := host2.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(t.Context()).Return([]host2.Host{}, nil)
			streamCmds := stream2.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(t.Context()).Return([]stream2.Stream{}, nil)
			cacheCmds := cache2.NewMockedCommands(ctrl)
			cacheCmds.EXPECT().GetAllInUse(t.Context()).Return([]cache2.Cache{}, nil)

			p1 := NewMockedfileProvider(ctrl)
			p1.EXPECT().provide(gomock.Any()).Return([]File{{Name: "f1.conf"}}, nil)
			p2 := NewMockedfileProvider(ctrl)
			p2.EXPECT().provide(gomock.Any()).Return([]File{{Name: "f2.conf"}}, nil)

			settingsCmds := settings2.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings2.Settings{}, nil).AnyTimes()

			facade := &Facade{
				hostCommands:     hostCmds,
				streamCommands:   streamCmds,
				cacheCommands:    cacheCmds,
				settingsCommands: settingsCmds,
				providers:        []fileProvider{p1, p2},
			}

			files, _, _, err := facade.GetConfigurationFiles(t.Context(), paths, features)
			assert.NoError(t, err)
			assert.Len(t, files, 2)
		})
	})

	t.Run("ReplaceConfigurationFiles", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tmpDir := t.TempDir()
		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.nginx.config-path": tmpDir,
		})

		features := &SupportedFeatures{
			StreamType: NoneSupportType,
		}

		t.Run("successfully replaces all files", func(t *testing.T) {
			hostCmds := host2.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(t.Context()).Return([]host2.Host{}, nil)
			streamCmds := stream2.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(t.Context()).Return([]stream2.Stream{}, nil)
			cacheCmds := cache2.NewMockedCommands(ctrl)
			cacheCmds.EXPECT().GetAllInUse(t.Context()).Return([]cache2.Cache{}, nil)

			provider := NewMockedfileProvider(ctrl)
			provider.EXPECT().
				provide(gomock.Any()).
				Return([]File{
					{
						Name:     "nginx.conf",
						Contents: "events {}",
					},
				}, nil)

			settingsCmds := settings2.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings2.Settings{}, nil).AnyTimes()

			facade := &Facade{
				hostCommands:     hostCmds,
				streamCommands:   streamCmds,
				cacheCommands:    cacheCmds,
				settingsCommands: settingsCmds,
				configuration:    cfg,
				providers:        []fileProvider{provider},
			}

			hosts, streams, err := facade.ReplaceConfigurationFiles(t.Context(), features)

			assert.NoError(t, err)
			assert.Empty(t, hosts)
			assert.Empty(t, streams)

			content, err := os.ReadFile(filepath.Join(tmpDir, "config", "nginx.conf"))
			require.NoError(t, err)
			assert.Equal(t, "events {}", string(content))

			assert.DirExists(t, filepath.Join(tmpDir, "logs"))
			assert.DirExists(t, filepath.Join(tmpDir, "cache"))
		})
	})
}
