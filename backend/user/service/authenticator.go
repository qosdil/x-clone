package service

type authenticator interface {
	GenerateToken(string) (string, error)
}
