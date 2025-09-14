package letsencrypt

import (
	"crypto"
	"crypto/rsa"

	"github.com/go-acme/lego/v4/registration"
)

type userDetails struct {
	email        string
	registration *registration.Resource
	privateKey   *rsa.PrivateKey
	newAccount   bool
}

func (u *userDetails) GetEmail() string {
	return u.email
}

func (u *userDetails) GetRegistration() *registration.Resource {
	return u.registration
}

func (u *userDetails) GetPrivateKey() crypto.PrivateKey {
	return u.privateKey
}
