package controller

import (
	"encoding/json"
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/utils"
)

type AuthenticationController struct {
	authenticationService service.AuthService
}

func NewAuthenticationController(authenticationService service.AuthService) *AuthenticationController {
	return &AuthenticationController{
		authenticationService: authenticationService,
	}
}

func (c *AuthenticationController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /authentication/sign-in", c.SignIn)
	mux.HandleFunc("POST /authentication/sign-up", c.SignUp)
	mux.HandleFunc("POST /authentication/refresh-token", c.RefreshToken)
}

// SignIn handles user login
func (c *AuthenticationController) SignIn(w http.ResponseWriter, r *http.Request) {
	var req model.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := c.authenticationService.SignIn(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	utils.WriteJSON(w, resp, http.StatusOK)
}

// SignUp handles user registration
func (c *AuthenticationController) SignUp(w http.ResponseWriter, req *http.Request) {
	var userSignUp model.SignUpRequest
	if err := json.NewDecoder(req.Body).Decode(&userSignUp); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !utils.IsEmail(userSignUp.Email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	if err := c.authenticationService.SignUp(userSignUp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *AuthenticationController) RefreshToken(w http.ResponseWriter, req *http.Request) {
	claims, err := c.authenticationService.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "invalid token", 400)
		return
	}
	isValid := true
	if !isValid {
		http.Error(w, "Invalid Token", 400)
		return
	}

	tokens, err := c.authenticationService.GenerateNewTokens(claims.UserId)
	if err != nil {
		http.Error(w, "could not create token", 500)
		return
	}
	utils.WriteJSON(w, map[string]string{"accessToken": tokens.AccessToken, "refreshToken": tokens.RefreshToken}, http.StatusOK)
}
