package jwt

import (
	"time"

	gwt "github.com/golang-jwt/jwt/v5"
)

const Issuer = "Skypiea AI"

// Generate
// validate
// invalidate
// parsefromRequest
// parsefromRequestWithClaims
// parseUser

type Manager interface {
	GenerateToken(userID, role string) (string, error)
	ValidateToken(token string) (*Claims, error)
	StoreToken(userID, token string, expiresIn time.Duration) error
	RevokeToken(token string) error
	IsTokenRevoked(token string) bool
}

type Config struct {
	Secret []byte
}

type Service struct {
	secret []byte
	store  StoreToken
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	gwt.RegisteredClaims
}

func NewJWTService(config *Config) *Service {
	if config == nil {
		panic("jwt config is nil")
	}
	return &Service{secret: config.Secret}
}

func (s *Service) GenerateToken(userID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: gwt.RegisteredClaims{
			Issuer:    Issuer,
			ExpiresAt: gwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  gwt.NewNumericDate(time.Now()),
		},
	}
	token := gwt.NewWithClaims(gwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := gwt.ParseWithClaims(tokenString, &Claims{}, func(token *gwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}
	return claims, nil
}

func (s *Service) StoreToken(userID, token string, expiresIn time.Duration) error {
	if s.store == nil {
		return nil
	}
	return s.store.StoreToken(userID, token, expiresIn)
}

func (s *Service) RevokeToken(token string) error {
	if s.store == nil {
		return nil
	}
	return s.store.RevokeToken(token)
}

func (s *Service) IsTokenRevoked(token string) bool {
	if s.store == nil {
		return false
	}
	return s.store.IsTokenRevoked(token)
}
