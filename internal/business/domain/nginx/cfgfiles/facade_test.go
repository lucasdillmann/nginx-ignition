package cfgfiles

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/nginx/cfgfiles/provider"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/stream"
)

func Test_Facade(t *testing.T) {
	t.Run("GetConfigurationFiles", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		paths := &provider.Paths{
			Base: "/etc/nginx/",
		}
		features := &provider.SupportedFeatures{
			StreamType: provider.StaticSupportType,
		}

		t.Run("successfully collects files from providers", func(t *testing.T) {
			hostCmds := host.NewMockedCommands(ctrl)
			hostCmds.EXPECT().
				GetAllEnabled(t.Context()).
				Return([]host.Host{
					{
						ID:          uuid.New(),
						DomainNames: []string{"example.com"},
					},
				}, nil)
			streamCmds := stream.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(t.Context()).Return([]stream.Stream{}, nil)
			cacheCmds := cache.NewMockedCommands(ctrl)
			cacheCmds.EXPECT().GetAllInUse(t.Context()).Return([]cache.Cache{}, nil)

			mockProvider := provider.NewMockedProvider(ctrl)
			mockProvider.EXPECT().
				Provide(gomock.Any()).
				Return([]provider.File{
					{
						Name:     "test.conf",
						Contents: "test",
					},
				}, nil)

			settingsCmds := settings.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(t.Context()).Return(&settings.Settings{}, nil)

			facade := &Facade{
				hostCommands:     hostCmds,
				streamCommands:   streamCmds,
				cacheCommands:    cacheCmds,
				settingsCommands: settingsCmds,
				providers:        []provider.Provider{mockProvider},
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
			hostCmds := host.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(t.Context()).Return(nil, assert.AnError)
			facade := &Facade{hostCommands: hostCmds}
			_, _, _, err := facade.GetConfigurationFiles(t.Context(), paths, features)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("returns error when settingsCommands fails", func(t *testing.T) {
			hostCmds := host.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(gomock.Any()).Return([]host.Host{}, nil)
			streamCmds := stream.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(gomock.Any()).Return([]stream.Stream{}, nil)
			cacheCmds := cache.NewMockedCommands(ctrl)
			cacheCmds.EXPECT().GetAllInUse(gomock.Any()).Return([]cache.Cache{}, nil)
			settingsCmds := settings.NewMockedCommands(ctrl)
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
			hostCmds := host.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(t.Context()).Return([]host.Host{}, nil)
			streamCmds := stream.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(t.Context()).Return(nil, assert.AnError)
			settingsCmds := settings.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings.Settings{}, nil).AnyTimes()
			facade := &Facade{
				hostCommands:     hostCmds,
				streamCommands:   streamCmds,
				settingsCommands: settingsCmds,
			}
			_, _, _, err := facade.GetConfigurationFiles(t.Context(), paths, features)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("returns error when cacheCommands fails", func(t *testing.T) {
			hostCmds := host.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(t.Context()).Return([]host.Host{}, nil)
			streamCmds := stream.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(t.Context()).Return([]stream.Stream{}, nil)
			cacheCmds := cache.NewMockedCommands(ctrl)
			cacheCmds.EXPECT().GetAllInUse(t.Context()).Return(nil, assert.AnError)
			settingsCmds := settings.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings.Settings{}, nil).AnyTimes()
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
			hostCmds := host.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(t.Context()).Return([]host.Host{}, nil)
			streamCmds := stream.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(t.Context()).Return([]stream.Stream{}, nil)
			cacheCmds := cache.NewMockedCommands(ctrl)
			cacheCmds.EXPECT().GetAllInUse(t.Context()).Return([]cache.Cache{}, nil)

			mockProvider := provider.NewMockedProvider(ctrl)
			mockProvider.EXPECT().Provide(gomock.Any()).Return(nil, assert.AnError)

			settingsCmds := settings.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings.Settings{}, nil).AnyTimes()

			facade := &Facade{
				hostCommands:     hostCmds,
				streamCommands:   streamCmds,
				cacheCommands:    cacheCmds,
				settingsCommands: settingsCmds,
				providers:        []provider.Provider{mockProvider},
			}

			_, _, _, err := facade.GetConfigurationFiles(t.Context(), paths, features)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("collects files from multiple providers", func(t *testing.T) {
			hostCmds := host.NewMockedCommands(ctrl)
			hostCmds.EXPECT().GetAllEnabled(t.Context()).Return([]host.Host{}, nil)
			streamCmds := stream.NewMockedCommands(ctrl)
			streamCmds.EXPECT().GetAllEnabled(t.Context()).Return([]stream.Stream{}, nil)
			cacheCmds := cache.NewMockedCommands(ctrl)
			cacheCmds.EXPECT().GetAllInUse(t.Context()).Return([]cache.Cache{}, nil)

			p1 := provider.NewMockedProvider(ctrl)
			p1.EXPECT().Provide(gomock.Any()).Return([]provider.File{{Name: "f1.conf"}}, nil)
			p2 := provider.NewMockedProvider(ctrl)
			p2.EXPECT().Provide(gomock.Any()).Return([]provider.File{{Name: "f2.conf"}}, nil)

			settingsCmds := settings.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings.Settings{}, nil).AnyTimes()

			facade := &Facade{
				hostCommands:     hostCmds,
				streamCommands:   streamCmds,
				cacheCommands:    cacheCmds,
				settingsCommands: settingsCmds,
				providers:        []provider.Provider{p1, p2},
			}

			files, _, _, err := facade.GetConfigurationFiles(t.Context(), paths, features)
			assert.NoError(t, err)
			assert.Len(t, files, 2)
		})
	})
}
