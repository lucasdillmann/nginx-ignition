package certificate

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type certificateModel struct {
	bun.BaseModel `bun:"certificate"`

	ID                 uuid.UUID  `bun:"id,pk"`
	DomainNames        []string   `bun:"domain_names,array"`
	ProviderID         string     `bun:"provider_id"`
	IssuedAt           time.Time  `bun:"issued_at"`
	ValidUntil         time.Time  `bun:"valid_until"`
	ValidFrom          time.Time  `bun:"valid_from"`
	RenewAfter         *time.Time `bun:"renew_after"`
	PrivateKey         string     `bun:"private_key"`
	PublicKey          string     `bun:"public_key"`
	CertificationChain []string   `bun:"certification_chain,array"`
	Parameters         string     `bun:"parameters"`
	Metadata           *string    `bun:"metadata"`
}
