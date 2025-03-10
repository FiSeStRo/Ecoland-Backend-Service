package config

import (
	"os"
	"strconv"
)

// JWTConfig holds all JWT-related configuration
type JWTConfig struct {
	Secret           []byte
	AccessTokenTime  int // in seconds
	RefreshTokenTime int // in seconds
	Issuer           string
}

// LoadJWTConfig loads JWT configuration from environment variables
func LoadJWTConfig() JWTConfig {
	const defaultAccessTime = 300
	const defaultRefreshTime = 3600
	accessTime, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_TIME_AT"))
	if err != nil {
		accessTime = defaultAccessTime
	}

	refreshTime, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_TIME_RT"))
	if err != nil {
		refreshTime = defaultRefreshTime
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-dev-secret"
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "EconyLand"
	}

	return JWTConfig{
		Secret:           []byte(secret),
		AccessTokenTime:  accessTime,
		RefreshTokenTime: refreshTime,
		Issuer:           issuer,
	}
}
