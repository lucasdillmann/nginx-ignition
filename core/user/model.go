package user

import "github.com/google/uuid"

type SaveRequest struct {
	ID          uuid.UUID
	Enabled     bool
	Name        string
	Username    string
	Password    *string
	Permissions Permissions
}

type AccessLevel string

const (
	NoAccessAccessLevel  AccessLevel = "NO_ACCESS"
	ReadOnlyAccessLevel  AccessLevel = "READ_ONLY"
	ReadWriteAccessLevel AccessLevel = "READ_WRITE"
)

type User struct {
	ID           uuid.UUID
	Enabled      bool
	Name         string
	Username     string
	PasswordHash string
	PasswordSalt string
	Permissions  Permissions
}

type Permissions struct {
	Hosts        AccessLevel
	Streams      AccessLevel
	Certificates AccessLevel
	Logs         AccessLevel
	Integrations AccessLevel
	AccessLists  AccessLevel
	Settings     AccessLevel
	Users        AccessLevel
	NginxServer  AccessLevel
}
