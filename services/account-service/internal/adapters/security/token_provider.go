package security

import (
	"account-service/internal/core/ports"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`

	jwt.RegisteredClaims
}

type JWTTokenProvider struct {
	jwtSecret     []byte
	tokenDuration time.Duration
}

func NewJWTTokenProvider(jwtSecret []byte, tokenDuration time.Duration) *JWTTokenProvider {
	return &JWTTokenProvider{
		jwtSecret:     jwtSecret,
		tokenDuration: tokenDuration,
	}
}

func (p *JWTTokenProvider) GenerateToken(userID, email string) (string, error) {

	claims := claims{
		UserID: userID,
		Email:  email,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(p.tokenDuration),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		p.jwtSecret,
	)

	if err != nil {
		return "", fmt.Errorf("error generating token: %v", err)
	}

	return tokenString, nil
}

func (p *JWTTokenProvider) ParseToken(tokenString string) (*ports.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return p.jwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &ports.TokenPayload{
		UserID: c.UserID,
		Email:  c.Email,
	}, nil
}

func (p *JWTTokenProvider) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (p *JWTTokenProvider) GenerateRefreshTokenWithHash() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (p *JWTTokenProvider) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (p *JWTTokenProvider) CompareToken(rawToken, hashedToken string) bool {
	computedHash := p.HashToken(rawToken)

	return subtle.ConstantTimeCompare([]byte(computedHash), []byte(hashedToken)) == 1
}
