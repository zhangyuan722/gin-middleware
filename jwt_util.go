package gm

import (
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

const tokenPrefix = "Bearer "

type Claims = jwt.RegisteredClaims
type NumericDate = jwt.NumericDate

func NewNumericDate(t time.Time) *NumericDate {
	return jwt.NewNumericDate(t)
}

func GenerateToken(claims Claims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))

	return tokenPrefix + tokenString, err
}

func ParseToken(tokenString string, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(strings.TrimPrefix(tokenString, tokenPrefix), &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if token != nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, err
}
