package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FiSeStRo/Ecoland-Backend-Service/authentication"
	"github.com/FiSeStRo/Ecoland-Backend-Service/database"
	"github.com/FiSeStRo/Ecoland-Backend-Service/utils"
)

func ListOfProductions(w http.ResponseWriter, req *http.Request) {

	if !utils.IsMethodPOST(w, req) {
		return
	}

	claims, err := authentication.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "invaldi token", http.StatusUnauthorized)
		return
	}

	var buildingId struct {
		BuildingId int `json:"buildingId"`
	}
	err = json.NewDecoder(req.Body).Decode(&buildingId)
	if err != nil {
		http.Error(w, "unable to parse JSON request", 400)
		return
	}

	db := database.GetDB()
	//find building
	row := db.QueryRow(`SELECT def_id FROM buildings WHERE id = ? AND user_id = ?)`, buildingId, claims.UserId)
	var buildingDefId int
	err = row.Scan(&buildingDefId)
	if err != nil {
		http.Error(w, "invalid building", http.StatusBadRequest)
	}

	rows, err := db.Query(`
				SELECT rel.production_id, p.token_name, p.cost, p.base_duration 
				FROM def_rel_building_production rel
				JOIN def_production p ON rel.production_id
				WHERE rel.building_id = ?
				`, buildingDefId)
	if err != nil {
		http.Error(w, "invalid production", http.StatusBadRequest)
		return
	}
	stmt, err := db.Prepare(`SELECT p.product_id, p.is_input, p.amount, def.token_name 
			FROM def_rel_production_product p 
			JOIN def_product def ON p.prodcut_id = def.id 
			WHERE p.production_id=?`)
	if err != nil {
		http.Error(w, "invalid sql", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	type BaseProduct struct {
		Id        int    `json:"id"`
		TokenName string `json:"token_name"`
		Amount    int    `json:"amount"`
	}
	type ProductionRes struct {
		Id           int           `json:"id"`
		TokenName    string        `json:"token_name"`
		BaseDuration time.Time     `json:"base_duration"`
		Cost         float64       `json:"cost"`
		ProductsIn   []BaseProduct `json:"products_in"`
		ProductsOut  []BaseProduct `json:"products_out"`
	}
	productionList := make([]ProductionRes, 0)
	for rows.Next() {
		var production ProductionRes
		err = rows.Scan(&production.Id, &production.TokenName, &production.Cost, &production.BaseDuration)
		if err != nil {
			http.Error(w, "error scanning productions", http.StatusInternalServerError)
			return
		}
		rslt, err := stmt.Query(production.Id)
		if err != nil {
			http.Error(w, "error executing production table call", http.StatusInternalServerError)
			return
		}

		for rslt.Next() {
			var baseProduct BaseProduct
			var isInput bool
			rslt.Scan(&baseProduct.Id, &isInput, &baseProduct.Amount, &baseProduct.TokenName)
			if isInput {
				production.ProductsIn = append(production.ProductsIn, baseProduct)
			} else {
				production.ProductsOut = append(production.ProductsOut, baseProduct)
			}
		}
		productionList = append(productionList, production)
	}
	utils.SetHeaderJson(w)
	json.NewEncoder(w).Encode(productionList)
}
