package authorization

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/apierror"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/user"
)

const (
	uniqueIdentifier           = "nginx-ignition"
	expectedJwtSecretSizeChars = 64
	expectedJwtSecretSizeBytes = 512
)

type Jwt struct {
	configuration *configuration.Configuration
	commands      user.Commands
	revokedIDs    []string
	secretKey     []byte
}

func newJwt(cfg *configuration.Configuration, commands user.Commands) (*Jwt, error) {
	prefixedConfiguration := cfg.WithPrefix("nginx-ignition.security.jwt")

	secretKey, err := initializeSecret(prefixedConfiguration)
	if err != nil {
		return nil, err
	}

	return &Jwt{
		configuration: prefixedConfiguration,
		commands:      commands,
		secretKey:     secretKey,
		revokedIDs:    []string{},
	}, nil
}

func (j *Jwt) RevokeToken(tokenID string) {
	j.revokedIDs = append(j.revokedIDs, tokenID)
}

func (j *Jwt) GenerateToken(usr *user.User) (*string, error) {
	ttlSeconds, err := j.configuration.GetInt("ttl-seconds")
	if err != nil {
		return nil, err
	}

	clockSkewSeconds, err := j.configuration.GetInt("clock-skew-seconds")
	if err != nil {
		return nil, err
	}

	notBefore := time.Now().Add(time.Second * time.Duration(clockSkewSeconds) * -1).Unix()
	expiresAt := time.Now().
		Add(time.Second * time.Duration(ttlSeconds)).
		Add(time.Second * time.Duration(clockSkewSeconds)).
		Unix()

	claims := jwt.MapClaims{
		"aud": uniqueIdentifier,
		"iss": uniqueIdentifier,
		"nbf": notBefore,
		"iat": time.Now().Unix(),
		"exp": expiresAt,
		"jti": uuid.New().String(),
		"sub": usr.ID.String(),
	}

	return j.sign(&claims)
}

func (j *Jwt) ValidateToken(ctx context.Context, tokenString string) (*Subject, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apierror.New(
				http.StatusUnauthorized,
				i18n.M(ctx, i18n.K.ApiCommonAuthorizationInvalidAccessToken),
			)
		}

		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["sub"].(string)
		userID, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		usr, err := j.commands.Get(ctx, userID)
		if err != nil {
			return nil, err
		}

		tokenID := claims["jti"].(string)
		if !usr.Enabled {
			j.RevokeToken(tokenID)
			return nil, apierror.New(
				http.StatusUnauthorized,
				i18n.M(ctx, i18n.K.ApiCommonAuthorizationInvalidAccessToken),
			)
		}

		if j.isRevoked(tokenID) {
			return nil, apierror.New(
				http.StatusUnauthorized,
				i18n.M(ctx, i18n.K.ApiCommonAuthorizationInvalidAccessToken),
			)
		}

		return &Subject{
			TokenID: tokenID,
			User:    usr,
			claims:  &claims,
		}, nil
	}

	return nil, apierror.New(
		http.StatusUnauthorized,
		i18n.M(ctx, i18n.K.ApiCommonAuthorizationInvalidAccessToken),
	)
}

func (j *Jwt) RefreshToken(subject *Subject) (*string, error) {
	windowSize, err := j.configuration.GetInt("refresh-window-seconds")
	if err != nil {
		return nil, err
	}

	clockSkewSeconds, err := j.configuration.GetInt("clock-skew-seconds")
	if err != nil {
		return nil, err
	}

	expiration, err := subject.claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	if time.Now().Add(time.Second * time.Duration(windowSize)).After(expiration.Time) {
		newClaims := *subject.claims
		newClaims["exp"] = time.Now().
			Add(time.Second * time.Duration(windowSize)).
			Add(time.Second * time.Duration(clockSkewSeconds)).
			Unix()
		return j.sign(&newClaims)
	}

	return nil, nil
}

func (j *Jwt) isRevoked(tokenID string) bool {
	for _, id := range j.revokedIDs {
		if id == tokenID {
			return true
		}
	}

	return false
}

func (j *Jwt) sign(claims *jwt.MapClaims) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(j.secretKey)
	return &tokenString, err
}

func initializeSecret(configurationProvider *configuration.Configuration) ([]byte, error) {
	secret, err := configurationProvider.Get("secret")
	if err != nil {
		secret = ""
	}

	if secret != "" {
		if len(secret) != expectedJwtSecretSizeChars {
			message := fmt.Sprintf(
				"JWT secret should be 64 characters long (512 bytes) but is %d characters long",
				len(secret),
			)
			return nil, errors.New(message)
		}

		return []byte(secret), nil
	}

	log.Warnf(
		"Application was initialized without a JWT secret and a random one will be generated. This will lead " +
			"to users being logged-out every time the app restarts or they hit a different instance. Please " +
			"refer to the documentation in order to provide a custom secret.",
	)

	secretBytes := make([]byte, expectedJwtSecretSizeBytes)
	_, err = rand.Read(secretBytes)
	if err != nil {
		return nil, err
	}

	return secretBytes, nil
}
