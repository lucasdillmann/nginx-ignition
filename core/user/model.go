package user

import "github.com/google/uuid"

type SaveRequest struct {
	Password    *string
	Permissions Permissions
	Name        string
	Username    string
	ID          uuid.UUID
	Enabled     bool
}

type AccessLevel string

const (
	NoAccessAccessLevel  AccessLevel = "NO_ACCESS"
	ReadOnlyAccessLevel  AccessLevel = "READ_ONLY"
	ReadWriteAccessLevel AccessLevel = "READ_WRITE"
)

type User struct {
	Permissions  Permissions
	Name         string
	Username     string
	PasswordHash string
	PasswordSalt string
	TOTP         TOTP
	ID           uuid.UUID
	Enabled      bool
}

type TOTP struct {
	Secret    *string
	Validated bool
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
	ExportData   AccessLevel
	VPNs         AccessLevel
	Caches       AccessLevel
	TrafficStats AccessLevel
}

type AuthenticationOutcome string

const (
	AuthenticationSuccessful  AuthenticationOutcome = "SUCCESS"
	AuthenticationFailed      AuthenticationOutcome = "FAILURE"
	AuthenticationMissingTOTP AuthenticationOutcome = "MISSING_TOTP"
)
