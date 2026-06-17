package security

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

// Claims e o payload que vai dentro do JWT: o id do usuario (no Subject padrao) e o papel dele
type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken assina um JWT novo (HS256) pro usuario, valido por ttl
func GenerateToken(secret string, userID int64, role string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken valida a assinatura/expiracao do token e devolve o id do usuario e o papel
func ParseToken(secret, tokenString string) (int64, string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, "", ErrInvalidToken
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, "", ErrInvalidToken
	}

	return userID, claims.Role, nil
}
