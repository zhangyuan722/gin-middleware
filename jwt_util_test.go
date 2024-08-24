package gm

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	tokenString, err := GenerateToken(Claims{
		ExpiresAt: NewNumericDate(time.Now().Add(expiresIn * time.Second)),
		Issuer:    Issuer,
		ID:        Id,
	}, secretKey)

	if err != nil {
		t.Errorf("Expected error nil, Got: %v", err.Error())
	}

	if tokenString == "" {
		t.Errorf("Expected valid token string, Got: empty")
	}
}

func TestParseToken(t *testing.T) {
	tokenString, _ := GenerateToken(Claims{
		ExpiresAt: NewNumericDate(time.Now().Add(expiresIn * time.Second)),
		Issuer:    Issuer,
		ID:        Id,
	}, secretKey)

	claims, err := ParseToken(tokenString, secretKey)

	if err != nil {
		t.Errorf("Expected error nil, Got: %v", err.Error())
	}

	if claims.ID != Id {
		t.Errorf("Expected %s, Got: %s", Id, claims.ID)
	}
}
