package authorization

import (
	"crypto/rand"
	"dillmann.com.br/nginx-ignition/api/common/api_error"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/user"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

const (
	uniqueIdentifier           = "nginx-ignition"
	expectedJwtSecretSizeChars = 64
	expectedJwtSecretSizeBytes = 512
)

var (
	errInvalidToken = api_error.New(http.StatusUnauthorized, "Invalid access token")
)

type Jwt struct {
	configuration *configuration.Configuration
	repository    *user.Repository
	revokedIds    []string
	secretKey     []byte
}

func newJwt(configuration *configuration.Configuration, repository *user.Repository) (*Jwt, error) {
	prefixedConfiguration := configuration.WithPrefix("nginx-ignition.security.jwt")

	secretKey, err := initializeSecret(prefixedConfiguration)
	if err != nil {
		return nil, err
	}

	return &Jwt{
		configuration: prefixedConfiguration,
		repository:    repository,
		secretKey:     secretKey,
		revokedIds:    []string{},
	}, nil
}

func (j *Jwt) RevokeToken(tokenId string) {
	j.revokedIds = append(j.revokedIds, tokenId)
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
		"aud":      uniqueIdentifier,
		"iss":      uniqueIdentifier,
		"nbf":      notBefore,
		"iat":      time.Now().Unix(),
		"exp":      expiresAt,
		"jti":      uuid.New().String(),
		"sub":      usr.ID.String(),
		"username": usr.Username,
		"role":     usr.Role,
	}

	return j.sign(&claims)
}

func (j *Jwt) ValidateToken(tokenString string) (*Subject, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidToken
		}

		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["sub"].(string)
		userId, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		usr, err := (*j.repository).FindByID(userId)
		if err != nil {
			return nil, err
		}

		tokenId := claims["jti"].(string)
		if !usr.Enabled {
			j.RevokeToken(tokenId)
			return nil, errInvalidToken
		}

		return &Subject{
			TokenID: tokenId,
			User:    usr,
			claims:  &claims,
		}, nil
	}

	return nil, errInvalidToken
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

func (j *Jwt) isRevoked(tokenId string) bool {
	for _, id := range j.revokedIds {
		if id == tokenId {
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

	log.Println(
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
