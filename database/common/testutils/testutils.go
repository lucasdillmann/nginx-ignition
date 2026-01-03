package testutils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"dillmann.com.br/nginx-ignition/database/common/migrations"
)

type TestExecutor func(*testing.T, *database.Database)

func RunWithMockedDatabases(t *testing.T, executor TestExecutor) {
	t.Run("SQLite", func(t *testing.T) {
		runWithMockedSQLiteDatabase(t, executor)
	})

	t.Run("PostgreSQL", func(t *testing.T) {
		runWithMockedPostgresDatabase(t, executor)
	})
}

func runWithMockedSQLiteDatabase(t *testing.T, executor TestExecutor) {
	overrides := map[string]string{
		"nginx-ignition.database.driver":    "sqlite",
		"nginx-ignition.database.data-path": t.TempDir(),
		"nginx-ignition.database.migrations-path": getEnvOrDefault(
			"TEST_MIGRATIONS_PATH",
			resolveMigrationsPath(),
		),
	}

	runWithMockedDatabase(t, overrides, executor, nil, nil)
}

func runWithMockedPostgresDatabase(t *testing.T, executor TestExecutor) {
	schemaName := fmt.Sprintf("test_%s", strings.ReplaceAll(uuid.New().String(), "-", ""))

	overrides := map[string]string{
		"nginx-ignition.database.driver": "postgres",
		"nginx-ignition.database.host": getEnvOrDefault(
			"TEST_POSTGRES_HOST",
			"localhost",
		),
		"nginx-ignition.database.port": getEnvOrDefault("TEST_POSTGRES_PORT", "5432"),
		"nginx-ignition.database.username": getEnvOrDefault(
			"TEST_POSTGRES_USER",
			"postgres",
		),
		"nginx-ignition.database.password": getEnvOrDefault(
			"TEST_POSTGRES_PASSWORD",
			"postgres",
		),
		"nginx-ignition.database.name": getEnvOrDefault(
			"TEST_POSTGRES_DBNAME",
			"nginx_ignition_tests",
		),
		"nginx-ignition.database.schema":   schemaName,
		"nginx-ignition.database.ssl-mode": "disable",
		"nginx-ignition.database.migrations-path": getEnvOrDefault(
			"TEST_MIGRATIONS_PATH",
			resolveMigrationsPath(),
		),
	}

	beforeMigrate := func(db *database.Database) {
		_, _ = db.Unwrap().Exec(fmt.Sprintf("CREATE SCHEMA %s", schemaName))
	}

	cleanup := func(db *database.Database) {
		_, _ = db.Unwrap().Exec(fmt.Sprintf("DROP SCHEMA %s CASCADE", schemaName))
	}

	runWithMockedDatabase(t, overrides, executor, beforeMigrate, cleanup)
}

func runWithMockedDatabase(
	t *testing.T,
	overrides map[string]string,
	executor TestExecutor,
	beforeMigrate func(*database.Database),
	cleanup func(*database.Database),
) {
	cfg := configuration.NewWithOverrides(overrides)
	db := database.New(cfg)

	err := db.Init()
	require.NoError(t, err)
	defer db.Close()

	if beforeMigrate != nil {
		beforeMigrate(db)
	}

	mgr := migrations.New(db, cfg)
	err = mgr.Migrate()
	require.NoError(t, err)

	executor(t, db)

	if cleanup != nil {
		cleanup(db)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func resolveMigrationsPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Warn("Unable to resolve migrations path: Failed to get caller information")
		return ""
	}

	baseDir := filepath.Dir(filename)
	return filepath.Join(baseDir, "..", "migrations", "scripts")
}
