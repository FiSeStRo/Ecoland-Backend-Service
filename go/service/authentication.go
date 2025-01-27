package service

import (
	"encoding/json"
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/authentication"
	"github.com/FiSeStRo/Ecoland-Backend-Service/utils"
)

// NewRefreshToken is a service to create a new Refresh and AccessToken for the client
func NewRefreshToken(w http.ResponseWriter, req *http.Request) {

	if !utils.IsMethodPOST(w, req) {
		return
	}

	claims, err := authentication.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "invalid token", 400)
		return
	}

	isValid := true
	if !isValid {
		http.Error(w, "Invalid Token", 400)
		return
	}

	atToken, err := authentication.CreateNewJWT(claims.UserId, false)
	if err != nil {
		http.Error(w, "could not create token", 500)
		return
	}
	rfToken, err := authentication.CreateNewJWT(claims.UserId, true)
	if err != nil {
		http.Error(w, "could not create token", 500)
		return
	}

	utils.SetHeaderJson(w)
	json.NewEncoder(w).Encode(map[string]string{"accessToken": atToken, "refreshToken": rfToken})

}
