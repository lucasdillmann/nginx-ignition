package backup

import (
	"context"
	"fmt"
	"os"

	"github.com/JCoupalK/go-pgdump"

	"dillmann.com.br/nginx-ignition/core/backup"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

type repository struct {
	db     *database.Database
	config *configuration.Configuration
}

func New(
	db *database.Database,
	config *configuration.Configuration,
) backup.Repository {
	return &repository{
		db:     db,
		config: config.WithPrefix("nginx-ignition.database"),
	}
}

func (r *repository) Get(_ context.Context) (*backup.Backup, error) {
	driver, err := r.config.Get("driver")
	if err != nil {
		return nil, err
	}

	switch driver {
	case "sqlite":
		return r.getSqliteBackup()
	case "postgres":
		return r.getPostgresBackup()
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}

func (r *repository) getSqliteBackup() (*backup.Backup, error) {
	fileName := "nginx-ignition.db"

	folder, err := r.config.Get("data-path")
	if err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("%s/%s", folder, fileName)
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return &backup.Backup{
		FileName:    fileName,
		ContentType: "application/x-sqlite3",
		Contents:    contents,
	}, nil
}

func (r *repository) getPostgresBackup() (*backup.Backup, error) {
	tempFile, err := os.CreateTemp(os.TempDir(), "nginx-ignition.sql")
	if err != nil {
		return nil, err
	}

	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	dumper := pgdump.NewDumper(r.db.ConnectionString(), 1)
	if err := dumper.DumpDatabase(tempFile.Name(), &pgdump.TableOptions{}); err != nil {
		return nil, err
	}

	contents, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return nil, err
	}

	return &backup.Backup{
		FileName:    "nginx-ignition.sql",
		ContentType: "text/plain",
		Contents:    contents,
	}, nil
}
