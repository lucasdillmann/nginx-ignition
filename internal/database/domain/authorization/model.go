package authorization

import (
	"github.com/uptrace/bun"
)

type configurationModel struct {
	bun.BaseModel `bun:"authorization_configuration"`

	JwtSecret string `bun:"jwt_secret,notnull"`
}
