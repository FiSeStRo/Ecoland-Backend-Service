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
				JOIN def_production p ON rel.production_id = p.id
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

func StartProduction(w http.ResponseWriter, req *http.Request) {

	type ReqBody struct {
		Id         int `json:"id"`
		BuildingId int `json:"building_id"`
		Cycles     int `json:"cycles"`
	}

	if !utils.IsMethodPOST(w, req) {
		return
	}

	claims, err := authentication.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	var body ReqBody
	err = json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	//check if building and userId exists and the building can produce the product

	user, err := database.FindUserById(claims.UserId)

	building, err := database.FindBuilding(body.BuildingId)
	if err != nil {
		http.Error(w, "no builfing found with this id", http.StatusBadRequest)
		return
	}
	if building.UserId != user.Id {
		http.Error(w, "no user with this building", http.StatusBadRequest)
		return
	}

	productions, err := database.GetPossibleProductionsOfBuilding(body.BuildingId)
	if err != nil {
		http.Error(w, "couldn't find productions", http.StatusInternalServerError)
	}

	var prod database.DefProduction
	for i := 0; i < len(productions); i++ {
		if productions[i].Id == body.Id {
			prod = productions[i]
			break
		}
	}
	userResources, err := database.GetUserResources(user.Id)
	if err != nil {
		http.Error(w, "no resources found", http.StatusBadRequest)
	}

	if prod.Cost > userResources.Money {
		http.Error(w, "insufficent Resources", http.StatusBadRequest)
		return
	}

	mats, err := database.FindProductsByProductionId(prod.Id)

	//check if the building has the needed resources
	storage, err := database.GetBuildingStorage(building.Id)
	if err != nil {
		http.Error(w, "Insufficent storage", http.StatusBadRequest)
	}

	//check fpr storage and reduce storage if possible
	for _, v := range mats {
		v.Amount *= body.Cycles
		for _, k := range storage {
			if k.ProductId == v.Id {
				if v.Amount <= k.Amount {
					k.Amount -= v.Amount
				} else {
					http.Error(w, "insufficent Storage", http.StatusBadRequest)
					return
				}
			}
		}
	}
	//start production
	db := database.GetDB()

	duration := prod.BaseDuration * body.Cycles
	rslt, err := db.Exec(`INSERT INTO rel_building_def_production(building_id, production_id, time_start, time_end, cycles, is_completed) VALUES(?,?,?,?,?,?)`, building.Id, prod.Id, time.Now().Unix(), time.Now().Unix()+int64(duration), body.Cycles, false)
	if err != nil {
		http.Error(w, "could not start production", http.StatusInternalServerError)
		return
	}

	resId, err := rslt.LastInsertId()
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
	utils.SetHeaderJson(w)
	json.NewEncoder(w).Encode(struct {
		Id int `json:"id"`
	}{
		Id: int(resId),
	})
}

func CancelProduction(w http.ResponseWriter, req *http.Request) {
	if !utils.IsMethodDELET(w, req) {
		return
	}
	claims, err := authentication.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	type ReqBody struct {
		BuildingId int `json:"id"`
	}
	//getUrlParam
	id, err := utils.GetUrlParamId(
		utils.UrlParam{Url: req.URL.Path, Position: 3},
	)
	if err != nil {
		if err.Error() == "wrong path" {
			http.Error(w, "wrong path", http.StatusNotFound)
		} else {
			http.Error(w, "invalid id", http.StatusBadRequest)
		}
		return
	}

	prod, err := database.FindProductionById(id)
	if err != nil {
		http.Error(w, "error finding production", http.StatusBadRequest)
		return
	}

	//Production belongs to user
	building, err := database.FindBuilding(prod.BuildingId)
	if err != nil {
		http.Error(w, "no building for producitonId", http.StatusBadRequest)
		return
	}
	if building.UserId != claims.UserId {
		http.Error(w, "invalid Request", http.StatusBadRequest)
		return
	}

	if prod.IsCompleted == true {
		http.Error(w, "can not cancel finished production", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	rslt, err := db.Exec(`DELETE FROM ? WHERE id=?`, database.BuildingProductionTable, id)
	if err != nil {
		http.Error(w, "Could not cancel Production", http.StatusInternalServerError)
		return
	}
	rows, err := rslt.RowsAffected()
	if err != nil || rows == 0 {
		http.Error(w, "production not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)

}
