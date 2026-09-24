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

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/apierror"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/ttlcache"
	user2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

const (
	uniqueIdentifier           = "nginx-ignition"
	expectedJwtSecretSizeChars = 64
	expectedJwtSecretSizeBytes = 512
)

type Jwt struct {
	commands           user2.Commands
	revokedTokens      *ttlcache.Cache[string, bool]
	secretKey          []byte
	ttlSeconds         int
	clockSkewSeconds   int
	renewWindowSeconds int
}

func newJwt(cfg *configuration.Configuration, commands user2.Commands) (*Jwt, error) {
	prefixedConfiguration := cfg.WithPrefix("nginx-ignition.security.jwt")

	secretKey, err := initializeSecret(prefixedConfiguration)
	if err != nil {
		return nil, err
	}

	ttlSeconds, err := prefixedConfiguration.GetInt("ttl-seconds")
	if err != nil {
		return nil, err
	}

	if ttlSeconds < 30 {
		return nil, errors.New("ttl-seconds cannot be less than 30")
	}

	clockSkewSeconds, err := prefixedConfiguration.GetInt("clock-skew-seconds")
	if err != nil {
		return nil, err
	}

	if clockSkewSeconds < 0 {
		return nil, errors.New("clock-skew-seconds cannot be negative")
	}

	renewWindowSeconds, err := prefixedConfiguration.GetInt("renew-window-seconds")
	if err != nil {
		return nil, err
	}

	if renewWindowSeconds < 0 {
		return nil, errors.New("renew-window-seconds cannot be negative")
	}

	if renewWindowSeconds > ttlSeconds {
		return nil, errors.New("renew-window-seconds cannot be bigger than ttl-seconds")
	}

	cacheTTL := (time.Duration(ttlSeconds) + time.Duration(clockSkewSeconds) + 1) * time.Second
	revokedTokens, err := ttlcache.New[string, bool](cacheTTL)
	if err != nil {
		return nil, err
	}

	return &Jwt{
		commands:           commands,
		secretKey:          secretKey,
		ttlSeconds:         ttlSeconds,
		clockSkewSeconds:   clockSkewSeconds,
		renewWindowSeconds: renewWindowSeconds,
		revokedTokens:      revokedTokens,
	}, nil
}

func (j *Jwt) RevokeToken(tokenID string) {
	j.revokedTokens.Set(tokenID, true)
}

func (j *Jwt) GenerateToken(usr *user2.User) (*string, error) {
	notBefore := time.Now().Add(time.Second * time.Duration(j.clockSkewSeconds) * -1).Unix()
	expiresAt := time.Now().
		Add(time.Second * time.Duration(j.ttlSeconds)).
		Add(time.Second * time.Duration(j.clockSkewSeconds)).
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
	expiration, err := subject.claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	if time.Now().Add(time.Second * time.Duration(j.renewWindowSeconds)).After(expiration.Time) {
		newClaims := make(jwt.MapClaims, len(*subject.claims)+1)
		for key, value := range *subject.claims {
			newClaims[key] = value
		}

		newClaims["jti"] = uuid.New().String()
		newClaims["exp"] = time.Now().
			Add(time.Second * time.Duration(j.renewWindowSeconds)).
			Add(time.Second * time.Duration(j.clockSkewSeconds)).
			Unix()

		result, err := j.sign(&newClaims)
		if err != nil {
			return nil, err
		}

		if previousJti, casted := (*subject.claims)["jti"].(string); casted {
			j.revokedTokens.Set(previousJti, true)
		}

		return result, nil
	}

	return nil, nil
}

func (j *Jwt) isRevoked(tokenID string) bool {
	_, ok := j.revokedTokens.Get(tokenID)
	return ok
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
