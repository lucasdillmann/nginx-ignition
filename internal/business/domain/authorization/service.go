package authorization

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
)

type service struct {
	repository    Repository
	configuration *configuration.Configuration
	secret        string
}

func newService(repository Repository, cfg *configuration.Configuration) *service {
	return &service{
		repository:    repository,
		configuration: cfg,
	}
}

func (s *service) JwtSecret() string {
	return s.secret
}

func (s *service) initialize(ctx context.Context) error {
	stored, err := s.repository.FindJwtSecret(ctx)
	if err != nil {
		return err
	}

	if stored != nil {
		return s.useStoredSecret(*stored)
	}

	secret := s.configuredSecret()
	if secret == "" {
		log.Warnf(
			"A secure random JWT secret was generated and stored in the database. Please refer to the " +
				"documentation to learn how to provide your own secret.",
		)

		if secret, err = generateSecret(); err != nil {
			return err
		}
	} else if len(secret) != jwtSecretSizeChars {
		return fmt.Errorf(
			"a custom JWT secret should be %d characters long, but the provided value is %d characters long",
			jwtSecretSizeChars,
			len(secret),
		)
	}

	s.secret = secret
	return s.repository.SaveJwtSecret(ctx, &secret)
}

func (s *service) useStoredSecret(stored string) error {
	if len(stored) != jwtSecretSizeChars {
		return fmt.Errorf(
			"JWT secret stored in the database should be %d characters long "+
				"but is %d characters long",
			jwtSecretSizeChars,
			len(stored),
		)
	}

	if configured := s.configuredSecret(); configured != "" && configured != stored {
		log.Warnf(
			"A JWT secret was provided through the configuration, but a different one was already stored " +
				"in the database. The stored one is the one being used. Please refer " +
				"to the documentation to learn how to change it (must be done directly in the database).",
		)
	}

	s.secret = stored
	return nil
}

func (s *service) configuredSecret() string {
	secret, err := s.configuration.Get("nginx-ignition.security.jwt.secret")
	if err != nil {
		return ""
	}

	return secret
}

func generateSecret() (string, error) {
	random := make([]byte, jwtSecretSizeChars/2)

	if _, err := rand.Read(random); err != nil {
		return "", err
	}

	return hex.EncodeToString(random), nil
}
