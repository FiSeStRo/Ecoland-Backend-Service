package authentication

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type EcoUserClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

// TODO: implement rotating keys with a db
var jWTKey = []byte("testKeyEcoLand")

// CreateNewJWT creates a new jwt with a userId payload and set's the expiration times based on if the creation should be a refreshToken
func CreateNewJWT(userId int, isRT bool) (string, error) {
	expirationTime := 5 * time.Minute
	if isRT {
		expirationTime = 12 * time.Minute
	}
	claims := EcoUserClaims{
		userId,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	st, err := token.SignedString(jWTKey)
	if err != nil {
		log.Println("jwt error", err)
		return "", fmt.Errorf("could not create and sign jwt %w", err)
	}

	return st, nil
}

// ValidateJWT validates a given jwt and return the EcoUserClaims if valid
func ValidateJWT(t string) (EcoUserClaims, error) {
	log.Println("token", t)
	vt, err := jwt.ParseWithClaims(t, &EcoUserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jWTKey, nil
	})

	isValid := err == nil && vt.Valid

	if !isValid {
		return EcoUserClaims{}, fmt.Errorf("jwt not valid %w", err)
	}

	claims := vt.Claims.(*EcoUserClaims)

	return *claims, nil

}

// ValidateAuthentication is a function for Authenticating a request
func ValidateAuthentication(req *http.Request) (EcoUserClaims, error) {
	authHeader := req.Header.Get("Authorization")

	if len(authHeader) < 7 && authHeader[:7] != "Bearer " {
		log.Println("wrong header", authHeader)
		return EcoUserClaims{}, fmt.Errorf("wrong authentication header")
	}

	token := authHeader[7:]

	return ValidateJWT(token)
}
