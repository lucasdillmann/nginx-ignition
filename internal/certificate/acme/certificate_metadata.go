package acme

type certificateMetadata struct {
	UserMail       string `json:"userMail"`
	UserPrivateKey string `json:"userPrivateKey"`
	UserPublicKey  string `json:"userPublicKey"`
}
