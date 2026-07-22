package ports

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hash, password string) error
}

type TokenPayload struct {
	UserID string
	Email  string
}

type TokenProvider interface {
	GenerateToken(userID, email string) (string, error)
	ParseToken(tokenString string) (*TokenPayload, error)
	GenerateRefreshToken() (string, error)
	HashToken(token string) string
	CompareToken(rawToken, hashedToken string) bool
}
