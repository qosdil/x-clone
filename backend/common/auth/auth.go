package auth

import "golang.org/x/crypto/bcrypt"

type auth struct{}

func NewAuth() *auth {
	return &auth{}
}

func (a *auth) CompareHashAndPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}

	return nil
}
