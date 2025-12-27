package certificate

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type certificateModel struct {
	bun.BaseModel `bun:"certificate"`

	ValidFrom          time.Time  `bun:"valid_from"`
	IssuedAt           time.Time  `bun:"issued_at"`
	ValidUntil         time.Time  `bun:"valid_until"`
	Metadata           *string    `bun:"metadata"`
	RenewAfter         *time.Time `bun:"renew_after"`
	ProviderID         string     `bun:"provider_id"`
	PrivateKey         string     `bun:"private_key"`
	PublicKey          string     `bun:"public_key"`
	Parameters         string     `bun:"parameters"`
	DomainNames        []string   `bun:"domain_names,array"`
	CertificationChain []string   `bun:"certification_chain,array"`
	ID                 uuid.UUID  `bun:"id,pk"`
}
