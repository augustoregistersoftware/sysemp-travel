package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret   []byte
	tokenTTL time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
}

func NewService(secret string, tokenTTL time.Duration) *Service {
	return &Service{
		secret:   []byte(secret),
		tokenTTL: tokenTTL,
	}
}

func (s *Service) GenerateToken(subject string) (string, error) {
	now := time.Now().UTC()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) ValidateToken(token string) (*Claims, error) {
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("token invalido")
	}

	if claims.Subject == "" {
		return nil, errors.New("token sem usuario")
	}

	return claims, nil
}

func (s *Service) ExpiresInSeconds() int64 {
	return int64(s.tokenTTL.Seconds())
}
