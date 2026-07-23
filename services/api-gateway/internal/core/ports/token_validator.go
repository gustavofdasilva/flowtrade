package ports

type TokenValidator interface {
	Validate(token string) (*Claims, error)
}

type Claims struct {
	UserID string
	Email  string
	Role   string
}
