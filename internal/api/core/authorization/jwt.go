package authorization

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/apierror"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

type Jwt struct {
	userCommands          user.Commands
	authorizationCommands authorization.Commands
	ttlSeconds            int
	clockSkewSeconds      int
	renewWindowSeconds    int
}

func newJwt(
	cfg *configuration.Configuration,
	userCommands user.Commands,
	authorizationCommands authorization.Commands,
) (*Jwt, error) {
	prefixedConfiguration := cfg.WithPrefix("nginx-ignition.security.jwt")

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

	return &Jwt{
		userCommands:          userCommands,
		authorizationCommands: authorizationCommands,
		ttlSeconds:            ttlSeconds,
		clockSkewSeconds:      clockSkewSeconds,
		renewWindowSeconds:    renewWindowSeconds,
	}, nil
}

func (j *Jwt) secretKey() []byte {
	return []byte(j.authorizationCommands.JwtSecret())
}

func (j *Jwt) GenerateToken(
	ctx context.Context,
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
		sessionExpiration := time.Now().
			Add(time.Second * time.Duration(j.ttlSeconds)).
			Add(time.Second * time.Duration(j.clockSkewSeconds))
		expiresAt = sessionExpiration.Unix()

		session, err := j.userCommands.StartSession(ctx, usr.ID, sessionExpiration)
		if err != nil {
			return nil, err
		}

		identifier = session.ID
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
		if err != nil {
			return nil, err
		}
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

		return j.secretKey(), nil
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

	if kind == SessionKind {
		if err = j.validateSession(ctx, tokenID, userID); err != nil {
			return nil, err
		}
	}

	usr, err := j.userCommands.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if usr == nil || !usr.Enabled {
		return nil, j.invalidTokenError(ctx)
	}

	return usr, nil
}

func (j *Jwt) validateSession(
	ctx context.Context,
	sessionID string,
	userID uuid.UUID,
) error {
	id, err := uuid.Parse(sessionID)
	if err != nil {
		return j.invalidTokenError(ctx)
	}

	session, err := j.userCommands.FindSession(ctx, id)
	if err != nil {
		return err
	}

	if session == nil || session.UserID != userID || time.Now().After(session.ExpiresAt) {
		return j.invalidTokenError(ctx)
	}

	return nil
}

func (j *Jwt) validateAPIToken(ctx context.Context, userID uuid.UUID, tokenID string) error {
	id, err := uuid.Parse(tokenID)
	if err != nil {
		return j.invalidTokenError(ctx)
	}

	token, err := j.userCommands.FindAPIToken(ctx, userID, id)
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

func (j *Jwt) RefreshToken(ctx context.Context, subject *Subject) (*string, error) {
	if subject == nil || subject.Kind != SessionKind || subject.claims == nil {
		return nil, errors.New("token cannot be refreshed")
	}

	expiration, err := subject.claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	if time.Now().Add(time.Second * time.Duration(j.renewWindowSeconds)).After(expiration.Time) {
		oldSessionID, err := uuid.Parse(subject.TokenID)
		if err != nil {
			return nil, err
		}

		expiresAt := time.Now().
			Add(time.Second * time.Duration(j.ttlSeconds)).
			Add(time.Second * time.Duration(j.clockSkewSeconds))

		session, err := j.userCommands.RotateSession(
			ctx,
			oldSessionID,
			subject.User.ID,
			expiresAt,
		)
		if err != nil {
			return nil, err
		}

		newClaims := make(jwt.MapClaims, len(*subject.claims))
		for key, value := range *subject.claims {
			newClaims[key] = value
		}

		newClaims["jti"] = session.ID.String()
		newClaims["exp"] = expiresAt.Unix()

		subject.TokenID = session.ID.String()
		subject.claims = &newClaims

		return j.sign(&newClaims)
	}

	return nil, nil
}

func (j *Jwt) sign(claims *jwt.MapClaims) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(j.secretKey())

	return &tokenString, err
}
