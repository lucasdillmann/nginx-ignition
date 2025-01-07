package password_hash

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"reflect"
	"slices"
)

const passwordHashIterations = 1024

func Hash(password string) (string, string, error) {
	salt := make([]byte, 64)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	hashResult, err := hashValue([]byte(password), salt)
	if err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(hashResult), base64.StdEncoding.EncodeToString(salt), nil
}

func Verify(password, hash, salt string) (bool, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return false, err
	}

	generatedHash, err := hashValue([]byte(password), saltBytes)
	if err != nil {
		return false, err
	}

	hashBytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(generatedHash, hashBytes), nil
}

func hashValue(password []byte, salt []byte) ([]byte, error) {
	output := slices.Concat(password, salt)
	hash := sha512.New()

	for range passwordHashIterations {
		if _, err := hash.Write(output); err != nil {
			return nil, err
		}
	}

	return output, nil
}
