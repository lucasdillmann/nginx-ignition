package letsencrypt

import (
	"crypto"
	"crypto/rsa"

	"github.com/go-acme/lego/v5/acme"
)

type userDetails struct {
	registration *acme.ExtendedAccount
	privateKey   *rsa.PrivateKey
	email        string
	newAccount   bool
}

func (u *userDetails) GetEmail() string {
	return u.email
}

func (u *userDetails) GetRegistration() *acme.ExtendedAccount {
	return u.registration
}

func (u *userDetails) GetPrivateKey() crypto.Signer {
	return u.privateKey
}
