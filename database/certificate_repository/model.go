package certificate_repository

import (
	"github.com/google/uuid"
	"time"
)

type certificateModel struct {
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
