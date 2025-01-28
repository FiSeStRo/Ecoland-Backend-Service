package authentication

import (
	"testing"
)

const userId = 1
const isRT = true
const nIsRt = false

var validToken string

func TestCreatNewJWT(t *testing.T) {

	t.Run("create NewJWT", func(t *testing.T) {
		st, err := CreateNewJWT(userId, isRT)

		if err != nil {
			t.Errorf("CreateNewJWT(%v, %v) returned with Error %q", userId, isRT, err)
		}

		validToken = st
	})

}

func TestValidateJWT(t *testing.T) {

	t.Run("valid JWT", func(t *testing.T) {
		claims, err := ValidateJWT(validToken)
		if err != nil {
			t.Errorf("the jwt could not be validated correctly")
		}
		if claims.UserId != 1 {
			t.Errorf("the jwt could not be validated, wrong userid, got %v, expected %v", claims.UserId, userId)
		}

	})

	t.Run("invalid JWT", func(t *testing.T) {
		_, err := ValidateJWT("wrongToken")
		if err == nil {
			t.Errorf("invalid token verified")
		}
	})

}
