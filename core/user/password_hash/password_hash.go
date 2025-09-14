package password_hash

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"reflect"
	"slices"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
)

type PasswordHash struct {
	configuration *configuration.Configuration
}

func New(configuration *configuration.Configuration) *PasswordHash {
	return &PasswordHash{configuration}
}

func (h *PasswordHash) Hash(password string) (string, string, error) {
	saltSize, hashIterations, err := readConfigValues(h.configuration)
	if err != nil {
		return "", "", err
	}

	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	hashResult, err := h.hashValue([]byte(password), salt, hashIterations)
	if err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(hashResult), base64.StdEncoding.EncodeToString(salt), nil
}

func (h *PasswordHash) Verify(password, hash, salt string) (bool, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return false, err
	}

	_, hashIterations, err := readConfigValues(h.configuration)
	if err != nil {
		return false, err
	}

	generatedHash, err := h.hashValue([]byte(password), saltBytes, hashIterations)
	if err != nil {
		return false, err
	}

	hashBytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(generatedHash, hashBytes), nil
}

func (h *PasswordHash) hashValue(password []byte, salt []byte, hashIterations int) ([]byte, error) {
	output := slices.Concat(password, salt)
	hash := sha512.New()

	for range hashIterations {
		if _, err := hash.Write(output); err != nil {
			return nil, err
		}
		output = hash.Sum(nil)
		hash.Reset()
	}

	return output, nil
}

func readConfigValues(configuration *configuration.Configuration) (int, int, error) {
	prefixedConfiguration := configuration.WithPrefix("nginx-ignition.security.user-password-hashing")
	saltSize, err := prefixedConfiguration.GetInt("salt-size")
	if err != nil {
		return 0, 0, err
	}

	iterations, err := prefixedConfiguration.GetInt("iterations")
	if err != nil {
		return 0, 0, err
	}

	return saltSize, iterations, nil
}
