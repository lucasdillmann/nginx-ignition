package backup

import "nginx-ignition/internal/core/backup"

func newBackup() *backup.Backup {
	return &backup.Backup{
		FileName:    "backup.zip",
		ContentType: "application/zip",
		Contents:    []byte("backup data"),
	}
}
