package auth

import (
	"fmt"
)

type auth struct{}

func NewAuth() *auth {
	return &auth{}
}

func (a *auth) GenerateToken(key string) (token string, err error) {
	token, err = GenerateJWT(key)
	if err != nil {
		err = fmt.Errorf("failed to generate JWT token: %v", err)
		return
	}

	return
}
