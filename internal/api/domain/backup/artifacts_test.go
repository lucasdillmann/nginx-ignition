package backup

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/backup"
)

func newBackup() *backup.Backup {
	return &backup.Backup{
		FileName:    "backup.zip",
		ContentType: "application/zip",
		Contents:    []byte("backup data"),
	}
}
