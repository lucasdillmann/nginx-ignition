package letsencrypt

type certificateMetadata struct {
	UserMail              string `json:"userMail"`
	UserPrivateKey        string `json:"userPrivateKey"`
	UserPublicKey         string `json:"userPublicKey"`
	ProductionEnvironment bool   `json:"productionEnvironment"`
}
