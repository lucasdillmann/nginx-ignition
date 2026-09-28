package nginx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/nginx/cfgfiles/provider"
)

func Test_service_configFiles(t *testing.T) {
	tmpDir := t.TempDir()
	nginxService := &service{
		processManager: &processManager{
			configPath: tmpDir,
		},
	}

	t.Run("configPaths", func(t *testing.T) {
		t.Run("builds paths relative to the config path", func(t *testing.T) {
			basePath := filepath.ToSlash(tmpDir) + "/"

			paths := nginxService.configPaths()

			assert.Equal(t, basePath, paths.Base)
			assert.Equal(t, basePath+"config/", paths.Config)
			assert.Equal(t, basePath+"logs/", paths.Logs)
			assert.Equal(t, basePath+"cache/", paths.Cache)
			assert.Equal(t, basePath+"temp/", paths.Temp)
		})

		t.Run("cleans up the config path", func(t *testing.T) {
			cleanedService := &service{
				processManager: &processManager{
					configPath: tmpDir + "/./config/..",
				},
			}

			paths := cleanedService.configPaths()

			assert.Equal(t, filepath.ToSlash(tmpDir)+"/", paths.Base)
		})
	})

	t.Run("newConfigPaths", func(t *testing.T) {
		t.Run("builds paths relative to the given root", func(t *testing.T) {
			root := filepath.Join(tmpDir, "candidate")

			paths := newConfigPaths(root)

			basePath := filepath.ToSlash(root) + "/"
			assert.Equal(t, basePath, paths.Base)
			assert.Equal(t, basePath+"config/", paths.Config)
			assert.Equal(t, basePath+"logs/", paths.Logs)
			assert.Equal(t, basePath+"cache/", paths.Cache)
			assert.Equal(t, basePath+"temp/", paths.Temp)
		})

		t.Run("cleans up the given root", func(t *testing.T) {
			paths := newConfigPaths(tmpDir + "/./config/..")

			assert.Equal(t, filepath.ToSlash(tmpDir)+"/", paths.Base)
		})
	})

	t.Run("createTestConfigRoot", func(t *testing.T) {
		t.Run("creates a folder inside the configured temp folder", func(t *testing.T) {
			paths := nginxService.configPaths()
			require.NoError(t, os.MkdirAll(filepath.Clean(paths.Temp), os.ModePerm))

			root, err := nginxService.createTestConfigRoot(paths)

			require.NoError(t, err)
			assert.DirExists(t, root)
			assert.Equal(t, filepath.Clean(paths.Temp), filepath.Dir(root))
		})

		t.Run("creates a different folder on every call", func(t *testing.T) {
			paths := nginxService.configPaths()
			require.NoError(t, os.MkdirAll(filepath.Clean(paths.Temp), os.ModePerm))

			first, err := nginxService.createTestConfigRoot(paths)
			require.NoError(t, err)
			t.Cleanup(func() { nginxService.removeTestConfigRoot(first) })

			second, err := nginxService.createTestConfigRoot(paths)
			require.NoError(t, err)
			t.Cleanup(func() { nginxService.removeTestConfigRoot(second) })

			assert.NotEqual(t, first, second)
		})

		t.Run("returns error when the temp folder is not usable", func(t *testing.T) {
			missingService := &service{
				processManager: &processManager{
					configPath: filepath.Join(tmpDir, "missing"),
				},
			}

			root, err := missingService.createTestConfigRoot(missingService.configPaths())

			assert.Error(t, err)
			assert.Empty(t, root)
		})
	})

	t.Run("removeTestConfigRoot", func(t *testing.T) {
		t.Run("removes the folder and its contents", func(t *testing.T) {
			root, err := nginxService.createTestConfigRoot(nginxService.configPaths())
			require.NoError(t, err)
			require.NoError(t, os.MkdirAll(filepath.Join(root, "config"), os.ModePerm))
			require.NoError(t, os.WriteFile(
				filepath.Join(root, "config", "nginx.conf"),
				[]byte("events {}"),
				0o644,
			))

			nginxService.removeTestConfigRoot(root)

			assert.NoDirExists(t, root)
		})

		t.Run("does nothing if the folder does not exist", func(t *testing.T) {
			nginxService.removeTestConfigRoot(filepath.Join(tmpDir, "config-test-missing"))

			assert.NoDirExists(t, filepath.Join(tmpDir, "config-test-missing"))
		})
	})

	t.Run("createMissingFolders", func(t *testing.T) {
		t.Run("creates config, logs, cache and temp folders", func(t *testing.T) {
			paths := nginxService.configPaths()

			err := nginxService.createMissingFolders(paths)

			assert.NoError(t, err)
			assert.DirExists(t, filepath.Clean(paths.Config))
			assert.DirExists(t, filepath.Clean(paths.Logs))
			assert.DirExists(t, filepath.Clean(paths.Cache))
			assert.DirExists(t, filepath.Clean(paths.Temp))
		})

		t.Run("keeps existing folders untouched", func(t *testing.T) {
			paths := nginxService.configPaths()
			marker := filepath.Join(filepath.Clean(paths.Logs), "marker.log")
			require.NoError(t, os.MkdirAll(filepath.Clean(paths.Logs), os.ModePerm))
			require.NoError(t, os.WriteFile(marker, nil, 0o644))

			err := nginxService.createMissingFolders(paths)

			assert.NoError(t, err)
			assert.FileExists(t, marker)
		})
	})

	t.Run("emptyConfigFolder", func(t *testing.T) {
		t.Run("removes all files from the config folder", func(t *testing.T) {
			paths := nginxService.configPaths()
			require.NoError(t, os.MkdirAll(filepath.Clean(paths.Config), os.ModePerm))
			leftover := filepath.Join(filepath.Clean(paths.Config), "leftover.conf")
			require.NoError(t, os.WriteFile(leftover, []byte("events {}"), 0o644))

			err := nginxService.emptyConfigFolder(paths)

			assert.NoError(t, err)
			assert.DirExists(t, filepath.Clean(paths.Config))
			assert.NoFileExists(t, leftover)
		})

		t.Run("returns error when config folder can not be read", func(t *testing.T) {
			missingService := &service{
				processManager: &processManager{
					configPath: filepath.Join(tmpDir, "missing"),
				},
			}

			err := missingService.emptyConfigFolder(missingService.configPaths())

			assert.Error(t, err)
		})
	})

	t.Run("writeConfigFile", func(t *testing.T) {
		t.Run("writes the formatted contents to the config folder", func(t *testing.T) {
			paths := nginxService.configPaths()
			require.NoError(t, os.MkdirAll(filepath.Clean(paths.Config), os.ModePerm))
			file := provider.File{
				Name:     "nginx.conf",
				Contents: "events {}",
			}

			err := nginxService.writeConfigFile(paths, file)

			assert.NoError(t, err)
			content, err := os.ReadFile(
				filepath.Join(filepath.Clean(paths.Config), "nginx.conf"),
			)
			require.NoError(t, err)
			assert.Equal(t, "events {}", string(content))
		})

		t.Run("returns error when file can not be written", func(t *testing.T) {
			paths := &provider.Paths{
				Config: filepath.Join(tmpDir, "missing", "config") + string(os.PathSeparator),
			}
			file := provider.File{
				Name:     "nginx.conf",
				Contents: "events {}",
			}

			err := nginxService.writeConfigFile(paths, file)

			assert.Error(t, err)
		})
	})
}
