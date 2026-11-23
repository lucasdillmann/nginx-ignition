package authorization

const (
	uniqueIdentifier = "nginx-ignition"
	tokenKindClaim   = "jtk"
	sessionIDClaim   = "sid"
)

type TokenKind int

const (
	SessionKind TokenKind = iota
	APIKind
)
