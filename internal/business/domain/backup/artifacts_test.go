package backup

func newBackup() *Backup {
	return &Backup{
		FileName:    "backup.db",
		ContentType: "application/octet-stream",
		Contents:    []byte("test content"),
	}
}
