package authorization

const (
	uniqueIdentifier = "nginx-ignition"
	tokenKindClaim   = "jtk"
)

type TokenKind int

const (
	SessionKind TokenKind = iota
	APIKind
)
