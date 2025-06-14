package model

// SignUpRequest represents user registration data
type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// SignInRequest represents user login data
type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignInResponse represents login response with tokens
type SignInResponse struct {
	ID           int    `json:"id"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
