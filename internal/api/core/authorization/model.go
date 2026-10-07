package authorization

const (
	uniqueIdentifier           = "nginx-ignition"
	expectedJwtSecretSizeChars = 64
	expectedJwtSecretSizeBytes = 512
	tokenKindClaim             = "jtk"
)

type TokenKind int

const (
	SessionKind TokenKind = iota
	APIKind
)
