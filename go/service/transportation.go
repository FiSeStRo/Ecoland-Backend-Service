package service

import (
	"encoding/json"
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/authentication"
	"github.com/FiSeStRo/Ecoland-Backend-Service/database"
	"github.com/FiSeStRo/Ecoland-Backend-Service/utils"
)

func ShipItems(w http.ResponseWriter, req *http.Request) {
	if !utils.IsMethodPOST(w, req) {
		return
	}
	claims, err := authentication.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	type ReqBody struct {
		FromId    int `json:"building_from_id"`
		ToId      int `json:"building_to_id"`
		ProductId int `json:"product_is" `
		Amount    int `json:"amount"`
	}

	var body ReqBody
	err = json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	//check if the buildings belong to the user
	if !database.IsUserBuildingOwner(claims.UserId, body.FromId) || !database.IsUserBuildingOwner(claims.UserId, body.ToId) {
		http.Error(w, "user does not own buildings", http.StatusBadRequest)
		return
	}

	fStorage, err := database.GetBStorageOfP(body.FromId, body.ProductId)
	if err != nil {
		http.Error(w, "could not find from storage", http.StatusBadRequest)
		return
	}
	tStorage, err := database.GetBStorageOfP(body.ToId, body.ProductId)
	if err != nil {
		http.Error(w, "could not find from storage", http.StatusBadRequest)
		return
	}

	if body.Amount > fStorage.Amount {
		http.Error(w, "from insufficent resources", http.StatusBadRequest)
		return
	}
	fStorage.Amount -= body.Amount
	tStorage.Amount += body.Amount
	if tStorage.Amount > tStorage.Capacity {
		http.Error(w, "not enough storage space at to", http.StatusBadRequest)
		return
	}
	err = database.BatchUpdateBStorageOfP([]database.BStorage{fStorage, tStorage})
	if err != nil {
		http.Error(w, "could not ship products", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
