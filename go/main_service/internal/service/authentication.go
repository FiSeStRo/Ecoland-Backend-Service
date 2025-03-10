package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/config"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	ValidateJWT(t string) (EcoUserClaims, error)
	ValidateAuthentication(req *http.Request) (EcoUserClaims, error)
	GenerateNewTokens(userId int) (Tokens, error)
	SignIn(req model.SignInRequest) (*model.SignInResponse, error)
	SignUp(req model.SignUpRequest) error
}

type authService struct {
	userRepo  repository.UserRepository
	jwtConfig config.JWTConfig
}

type EcoUserClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// TODO: implement rotating keys with a redisdb

func NewAuthService(jwtConfig config.JWTConfig, userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

// CreateAccessToken generates a new access token for a user
func (s *authService) CreateAccessToken(userId int) (string, error) {
	expiry := time.Duration(s.jwtConfig.AccessTokenTime) * time.Second
	return s.createToken(userId, expiry)
}

// CreateRefreshToken generates a new RefreshToken for a user
func (s *authService) CreateRefreshToken(userId int) (string, error) {
	expiry := time.Duration(s.jwtConfig.RefreshTokenTime) * time.Second
	return s.createToken(userId, expiry)
}

// createToken is a helper method to create tokens with different expiry times
func (s *authService) createToken(userId int, expiry time.Duration) (string, error) {
	claims := EcoUserClaims{
		userId,
		jwt.RegisteredClaims{
			Issuer:    s.jwtConfig.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	st, err := token.SignedString(s.jwtConfig.Secret)
	if err != nil {
		log.Println("jwt error", err)
		return "", fmt.Errorf("could not create and sign jwt %w", err)
	}

	return st, nil

}

// genrateNewTokens creates new Tokens for the user
func (s *authService) GenerateNewTokens(userId int) (Tokens, error) {

	atToken, err := s.CreateAccessToken(userId)
	if err != nil {
		return Tokens{}, err
	}
	rfToken, err := s.CreateRefreshToken(userId)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  atToken,
		RefreshToken: rfToken,
	}, nil

}

// SignUp handles user registration
func (s *authService) SignUp(req model.SignUpRequest) error {
	exists, err := s.userRepo.CheckUserExists(req.Username)
	if err != nil {
		return fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return errors.New("username already exists")
	}
	exists, err = s.userRepo.CheckEmailExists(req.Email)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return errors.New("email already exists")
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user := &model.User{
		Username: req.Username,
		Password: hashedPw,
		Email:    req.Email,
	}
	userID, err := s.userRepo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	err = s.userRepo.CreateUserResources(userID, 100000, 10)
	if err != nil {
		// Ideally we'd have transaction support here to rollback the user creation
		return fmt.Errorf("failed to create user resources: %w", err)
	}

	log.Printf("User created: %s (ID: %d)\n", user.Username, userID)
	return nil
}

// SignIn handles user authentication
func (s *authService) SignIn(req model.SignInRequest) (*model.SignInResponse, error) {
	exists, err := s.userRepo.CheckUserExists(req.Username)
	if err != nil || !exists {
		return nil, fmt.Errorf("invalid credentials")
	}
	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	tokens, err := s.GenerateNewTokens(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &model.SignInResponse{
		ID:           user.ID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// ValidateJWT validates a given jwt and return the EcoUserClaims if valid
func (s *authService) ValidateJWT(t string) (EcoUserClaims, error) {
	vt, err := jwt.ParseWithClaims(t, &EcoUserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.jwtConfig.Secret, nil
	})

	isValid := err == nil && vt.Valid
	if !isValid {
		return EcoUserClaims{}, fmt.Errorf("jwt not valid %w", err)
	}

	claims := vt.Claims.(*EcoUserClaims)

	return *claims, nil

}

// ValidateAuthentication is a function for Authenticating a request
func (s *authService) ValidateAuthentication(req *http.Request) (EcoUserClaims, error) {
	authHeader := req.Header.Get("Authorization")

	if len(authHeader) < 7 && authHeader[:7] != "Bearer " {
		log.Println("wrong header format")
		return EcoUserClaims{}, fmt.Errorf("wrong authentication header")
	}

	token := authHeader[7:]

	return s.ValidateJWT(token)
}
