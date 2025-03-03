package jwt

import "time"

type StoreToken interface {
	StoreToken(userID, token string, expiresIn time.Duration) error
	RevokeToken(token string) error
	IsTokenRevoked(token string) bool
}
