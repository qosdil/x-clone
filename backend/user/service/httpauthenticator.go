package service

type httpauthenticator interface {
	GenerateToken(string) (string, error)
}
