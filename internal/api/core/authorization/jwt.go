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
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

type Jwt struct {
	commands           user.Commands
	revokedTokens      *ttlcache.Cache[string, bool]
	secretKey          []byte
	ttlSeconds         int
	clockSkewSeconds   int
	renewWindowSeconds int
}

func newJwt(cfg *configuration.Configuration, commands user.Commands) (*Jwt, error) {
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

func (j *Jwt) GenerateToken(
	usr *user.User,
	kind TokenKind,
	tokenID *uuid.UUID,
	expiration *time.Time,
) (*string, error) {
	identifier := uuid.New()
	if tokenID != nil {
		identifier = *tokenID
	}

	var expiresAt int64
	if kind == SessionKind {
		expiresAt = time.Now().
			Add(time.Second * time.Duration(j.ttlSeconds)).
			Add(time.Second * time.Duration(j.clockSkewSeconds)).
			Unix()
	} else if expiration != nil {
		expiresAt = expiration.Unix()
	}

	claims := jwt.MapClaims{
		"aud":          uniqueIdentifier,
		"iss":          uniqueIdentifier,
		"nbf":          time.Now().Add(time.Second * time.Duration(j.clockSkewSeconds) * -1).Unix(),
		"iat":          time.Now().Unix(),
		"jti":          identifier.String(),
		"sub":          usr.ID.String(),
		tokenKindClaim: kind,
	}

	if expiresAt > 0 {
		claims["exp"] = expiresAt
	}

	return j.sign(&claims)
}

func (j *Jwt) ValidateToken(ctx context.Context, tokenString string) (*Subject, error) {
	token, err := j.parse(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	claims, isMapClaims := token.Claims.(jwt.MapClaims)
	if !isMapClaims || !token.Valid {
		return nil, j.invalidTokenError(ctx)
	}

	kind, err := j.resolveKind(ctx, claims)
	if err != nil {
		return nil, err
	}

	tokenID, casted := claims["jti"].(string)
	if !casted || tokenID == "" {
		return nil, j.invalidTokenError(ctx)
	}

	usr, err := j.resolveUser(ctx, claims, kind, tokenID)
	if err != nil {
		return nil, err
	}

	if kind == APIKind {
		err = j.validateAPIToken(ctx, usr.ID, tokenID)
	} else if j.isRevoked(tokenID) {
		err = j.invalidTokenError(ctx)
	}

	if err != nil {
		return nil, err
	}

	return &Subject{
		Kind:    kind,
		TokenID: tokenID,
		User:    usr,
		claims:  &claims,
	}, nil
}

func (j *Jwt) parse(ctx context.Context, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, j.invalidTokenError(ctx)
		}

		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (j *Jwt) resolveKind(ctx context.Context, claims jwt.MapClaims) (TokenKind, error) {
	value, casted := claims[tokenKindClaim].(float64)
	if !casted {
		return 0, j.invalidTokenError(ctx)
	}

	switch kind := TokenKind(value); kind {
	case SessionKind, APIKind:
		return kind, nil
	default:
		return 0, j.invalidTokenError(ctx)
	}
}

func (j *Jwt) resolveUser(
	ctx context.Context,
	claims jwt.MapClaims,
	kind TokenKind,
	tokenID string,
) (*user.User, error) {
	id, casted := claims["sub"].(string)
	if !casted {
		return nil, j.invalidTokenError(ctx)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	usr, err := j.commands.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !usr.Enabled && kind == SessionKind {
		j.RevokeToken(tokenID)
	}

	if !usr.Enabled {
		return nil, j.invalidTokenError(ctx)
	}

	return usr, nil
}

func (j *Jwt) validateAPIToken(ctx context.Context, userID uuid.UUID, tokenID string) error {
	id, err := uuid.Parse(tokenID)
	if err != nil {
		return j.invalidTokenError(ctx)
	}

	token, err := j.commands.FindAPIToken(ctx, userID, id)
	if err != nil {
		return err
	}

	if token == nil {
		return j.invalidTokenError(ctx)
	}

	if token.Expiration != nil && !token.Expiration.After(time.Now()) {
		return j.invalidTokenError(ctx)
	}

	return nil
}

func (j *Jwt) invalidTokenError(ctx context.Context) error {
	return apierror.New(
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
