package jwt

import (
	"api-gateway/internal/core/ports"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTValidator struct {
	secret string
}

type claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`

	jwt.RegisteredClaims
}

func NewJWTValidator(secret string) *JWTValidator {
	return &JWTValidator{
		secret: secret,
	}
}

func (v *JWTValidator) Validate(tokenString string) (*ports.Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(v.secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &ports.Claims{
		UserID: c.UserID,
		Email:  c.Email,
	}, nil
}
