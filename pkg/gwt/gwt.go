package gwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	Issuer   = "Skypiea AI Auth"
	CtxToken = "Token"
)

// TODO: currently we dont need interface here, caller side will be better fit
type Manager interface {
	GenerateToken(userID, role string) (string, error)
	ValidateToken(token string) (*Token, error)
}

type Config struct {
	Secret string `yaml:"jwt.secret"`
	Domain string `yaml:"jwt.domain"`
}

type Service struct {
	secret []byte
}

type Token struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(config *Config) *Service {
	if config == nil {
		panic("gwt config is nil")
	}
	return &Service{secret: []byte(config.Secret)}
}

func (s *Service) GenerateToken(userID, role string) (string, error) {
	claims := Token{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// TODO claims to token
func (s *Service) ValidateToken(tokenString string) (*Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Token)
	if !ok {
		return nil, err
	}
	return claims, nil
}
