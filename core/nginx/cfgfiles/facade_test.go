package cfgfiles

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/stream"
)

func Test_Facade_GetConfigurationFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	paths := &Paths{
		Base: "/etc/nginx/",
	}
	features := &SupportedFeatures{
		StreamType: StaticSupportType,
	}

	t.Run("successfully collects files from providers", func(t *testing.T) {
		hostCmds := &host.Commands{
			GetAllEnabled: func(_ context.Context) ([]host.Host, error) {
				return []host.Host{{ID: uuid.New(), DomainNames: []string{"example.com"}}}, nil
			},
		}
		streamCmds := &stream.Commands{
			GetAllEnabled: func(_ context.Context) ([]stream.Stream, error) {
				return []stream.Stream{}, nil
			},
		}
		cacheCmds := &cache.Commands{
			GetAllInUse: func(_ context.Context) ([]cache.Cache, error) {
				return []cache.Cache{}, nil
			},
		}

		provider := NewMockfileProvider(ctrl)
		provider.EXPECT().
			provide(gomock.Any()).
			Return([]File{{Name: "test.conf", Contents: "test"}}, nil)

		f := &Facade{
			hostCommands:   hostCmds,
			streamCommands: streamCmds,
			cacheCommands:  cacheCmds,
			providers:      []fileProvider{provider},
		}

		configFiles, hosts, streams, err := f.GetConfigurationFiles(ctx, paths, features)

		assert.NoError(t, err)
		assert.Len(t, configFiles, 1)
		assert.Equal(t, "test.conf", configFiles[0].Name)
		assert.Len(t, hosts, 1)
		assert.Len(t, streams, 0)
	})

	t.Run("returns error when hostCommands fails", func(t *testing.T) {
		hostCmds := &host.Commands{
			GetAllEnabled: func(_ context.Context) ([]host.Host, error) {
				return nil, assert.AnError
			},
		}
		f := &Facade{hostCommands: hostCmds}
		_, _, _, err := f.GetConfigurationFiles(ctx, paths, features)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("returns error when streamCommands fails", func(t *testing.T) {
		hostCmds := &host.Commands{
			GetAllEnabled: func(_ context.Context) ([]host.Host, error) { return []host.Host{}, nil },
		}
		streamCmds := &stream.Commands{
			GetAllEnabled: func(_ context.Context) ([]stream.Stream, error) {
				return nil, assert.AnError
			},
		}
		f := &Facade{
			hostCommands:   hostCmds,
			streamCommands: streamCmds,
		}
		_, _, _, err := f.GetConfigurationFiles(ctx, paths, features)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("returns error when cacheCommands fails", func(t *testing.T) {
		hostCmds := &host.Commands{
			GetAllEnabled: func(_ context.Context) ([]host.Host, error) { return []host.Host{}, nil },
		}
		streamCmds := &stream.Commands{
			GetAllEnabled: func(_ context.Context) ([]stream.Stream, error) { return []stream.Stream{}, nil },
		}
		cacheCmds := &cache.Commands{
			GetAllInUse: func(_ context.Context) ([]cache.Cache, error) {
				return nil, assert.AnError
			},
		}
		f := &Facade{
			hostCommands:   hostCmds,
			streamCommands: streamCmds,
			cacheCommands:  cacheCmds,
		}
		_, _, _, err := f.GetConfigurationFiles(ctx, paths, features)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("returns error when provider fails", func(t *testing.T) {
		hostCmds := &host.Commands{
			GetAllEnabled: func(_ context.Context) ([]host.Host, error) { return []host.Host{}, nil },
		}
		streamCmds := &stream.Commands{
			GetAllEnabled: func(_ context.Context) ([]stream.Stream, error) { return []stream.Stream{}, nil },
		}
		cacheCmds := &cache.Commands{
			GetAllInUse: func(_ context.Context) ([]cache.Cache, error) { return []cache.Cache{}, nil },
		}

		provider := NewMockfileProvider(ctrl)
		provider.EXPECT().provide(gomock.Any()).Return(nil, assert.AnError)

		f := &Facade{
			hostCommands:   hostCmds,
			streamCommands: streamCmds,
			cacheCommands:  cacheCmds,
			providers:      []fileProvider{provider},
		}

		_, _, _, err := f.GetConfigurationFiles(ctx, paths, features)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("collects files from multiple providers", func(t *testing.T) {
		hostCmds := &host.Commands{
			GetAllEnabled: func(_ context.Context) ([]host.Host, error) { return []host.Host{}, nil },
		}
		streamCmds := &stream.Commands{
			GetAllEnabled: func(_ context.Context) ([]stream.Stream, error) { return []stream.Stream{}, nil },
		}
		cacheCmds := &cache.Commands{
			GetAllInUse: func(_ context.Context) ([]cache.Cache, error) { return []cache.Cache{}, nil },
		}

		p1 := NewMockfileProvider(ctrl)
		p1.EXPECT().provide(gomock.Any()).Return([]File{{Name: "f1.conf"}}, nil)
		p2 := NewMockfileProvider(ctrl)
		p2.EXPECT().provide(gomock.Any()).Return([]File{{Name: "f2.conf"}}, nil)

		f := &Facade{
			hostCommands:   hostCmds,
			streamCommands: streamCmds,
			cacheCommands:  cacheCmds,
			providers:      []fileProvider{p1, p2},
		}

		files, _, _, err := f.GetConfigurationFiles(ctx, paths, features)
		assert.NoError(t, err)
		assert.Len(t, files, 2)
	})
}

func Test_Facade_ReplaceConfigurationFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	tmpDir := t.TempDir()
	t.Setenv("NGINX_IGNITION_NGINX_CONFIG_PATH", tmpDir)
	cfg := configuration.New()

	features := &SupportedFeatures{
		StreamType: NoneSupportType,
	}

	t.Run("successfully replaces all files", func(t *testing.T) {
		hostCmds := &host.Commands{
			GetAllEnabled: func(_ context.Context) ([]host.Host, error) { return []host.Host{}, nil },
		}
		streamCmds := &stream.Commands{
			GetAllEnabled: func(_ context.Context) ([]stream.Stream, error) { return []stream.Stream{}, nil },
		}
		cacheCmds := &cache.Commands{
			GetAllInUse: func(_ context.Context) ([]cache.Cache, error) { return []cache.Cache{}, nil },
		}

		provider := NewMockfileProvider(ctrl)
		provider.EXPECT().
			provide(gomock.Any()).
			Return([]File{{Name: "nginx.conf", Contents: "events {}"}}, nil)

		f := &Facade{
			hostCommands:   hostCmds,
			streamCommands: streamCmds,
			cacheCommands:  cacheCmds,
			configuration:  cfg,
			providers:      []fileProvider{provider},
		}

		hosts, streams, err := f.ReplaceConfigurationFiles(ctx, features)

		assert.NoError(t, err)
		assert.Empty(t, hosts)
		assert.Empty(t, streams)

		content, err := os.ReadFile(filepath.Join(tmpDir, "config", "nginx.conf"))
		require.NoError(t, err)
		assert.Equal(t, "events {}", string(content))

		assert.DirExists(t, filepath.Join(tmpDir, "logs"))
		assert.DirExists(t, filepath.Join(tmpDir, "cache"))
	})
}
