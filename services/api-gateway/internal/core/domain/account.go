package domain

type RegisterInput struct {
	Email    string
	Password string
	Username string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	AccessToken  string
	RefreshToken string
	User         UserOutput
}

type UserOutput struct {
	ID       string
	Username string
	Email    string
}

type UpdateUserInput struct {
	ID       string
	Username string
	Email    string
	Password string
}

type RefreshTokenInput struct {
	RefreshToken string
}
