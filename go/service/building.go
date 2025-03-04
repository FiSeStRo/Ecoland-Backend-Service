package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/FiSeStRo/Ecoland-Backend-Service/authentication"
	"github.com/FiSeStRo/Ecoland-Backend-Service/database"
	"github.com/FiSeStRo/Ecoland-Backend-Service/utils"
)

type BuildingBase struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type struct {
		Def_id     int `json:"def_id"`
		Token_name int `json:"token_name"`
	}
}

type Building struct {
	Building []BuildingBase
	Cost     float64 `json:"cost"`
	Time     int     `json:"time"`
}

type Def_Building struct {
	Id   int     `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
	Time int     `json:"time"`
}

type ConstructRequest struct {
	Def_id       int    `json:"def_id"`
	Display_name string `json:"display_name"`
	Lan          int    `json:"lan"`
	Lat          int    `json:"lat"`
}

// ConstrucitonList service to retrive a list of possible constructions
func ConstructionList(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Wrong method", 400)
	}
	db := database.GetDB()
	rows, err := db.Query(`SELECT id, token_name, base_construction_cost, base_construction_time from def_buildings`)
	if err != nil {
		http.Error(w, "Error reading from Buildings table", 400)
		return
	}
	defer rows.Close()

	var defBuildings []Def_Building

	for rows.Next() {
		var building Def_Building
		err := rows.Scan(
			&building.Id,
			&building.Name,
			&building.Cost,
			&building.Time,
		)

		if err != nil {
			http.Error(w, "Error scanning building row", 500)
			return
		}

		defBuildings = append(defBuildings, building)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(defBuildings)
}

// ConstructBuildings services to create a building for a user
func ConstructBuilding(w http.ResponseWriter, req *http.Request) {

	if !utils.IsMethodPOST(w, req) {
		return
	}

	claims, err := authentication.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
	}

	var constructReq ConstructRequest

	err = json.NewDecoder(req.Body).Decode(&constructReq)
	if err != nil {
		http.Error(w, "unable to parse JSON", 400)
		return
	}
	db := database.GetDB()
	//can user affort building
	var cost float64
	var money float64

	row := db.QueryRow(`SELECT base_consturction_cost FROM def_buildings WHERE id=?`, constructReq.Def_id)
	err = row.Scan(&cost)
	if err != nil {
		http.Error(w, "Error finding building cost", http.StatusInternalServerError)
		return
	}

	row = db.QueryRow(`SELECT money FROM user_resourcres WHERE id=?`, claims.UserId)
	err = row.Scan(&money)
	if err != nil {
		http.Error(w, "Error finding user money", http.StatusInternalServerError)
		return
	}

	if money < cost {
		http.Error(w, "insuuficent resources", http.StatusBadRequest)
		return
	}
	money -= cost

	err = database.UpdateUserResources(database.PUserResource{UserId: claims.UserId, Money: &money})
	if err != nil {
		http.Error(w, "could not update money", http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare(`INSERT INTO buildings (user_id, def_id, name, lan, lat time_build)
	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`)
	if err != nil {
		log.Println("Error preparing data insert", err)
		http.Error(w, "Error Preparing data insert", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(claims.UserId, constructReq.Def_id, constructReq.Display_name, constructReq.Lan, constructReq.Lat)
	if err != nil {
		http.Error(w, "Error constructing new building", http.StatusInternalServerError)
		return
	}
	buildingId, _ := result.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": buildingId})
}

type ProductInBuilding struct {
	Id        int    `json:"id"`
	TokenName string `json:"token_name"`
	Amount    struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"amount"`
}

type Production struct {
	Id     int `json:"id"`
	Cycles struct {
		Total     int `json:"total"`
		Completed int `json:"completed"`
	} `json:"cycles"`
	IsActive    bool                `json:"isActive"`
	TimeEnd     time.Time           `json:"time_end"`
	ProductsIn  []ProductInBuilding `json:"products_in"`
	ProductsOut []ProductInBuilding `json:"prducts_out"`
}

type BuildingListRes struct {
	BuildingBase
	TimeBuild  time.Time `json:"time_build"`
	Production Production
}

// ListOfBuidlings returns a list of possible buildings that the user can construct
func ListOfBuildings(w http.ResponseWriter, req *http.Request) {
	if !utils.IsMethodGET(w, req) {
		return
	}

	claims, err := authentication.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
	}

	db := database.GetDB()
	stmt, err := db.Prepare(`
		SELECT b.id, b.name, b.time_build, db.token_name, b.def_id 
        FROM buildings b
        JOIN def_buildings db ON b.def_id = db.id
        WHERE b.user_id = ?
		`)
	if err != nil {
		http.Error(w, "error preparing query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(claims.UserId)
	if err != nil {
		http.Error(w, "error searching for buildings", http.StatusInternalServerError)
		return
	}

	buildingList := make([]BuildingListRes, 0)
	for rows.Next() {
		b := BuildingListRes{}
		err = rows.Scan(&b.Id, &b.Name, &b.TimeBuild, &b.Type.Token_name, &b.Type.Def_id)
		if err != nil {
			log.Println("error", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		//get current Production of the building
		var p struct {
			id          int
			timeStart   time.Time
			timeEnd     time.Time
			cycles      int
			isCompleted bool
		}
		rowp := db.QueryRow(`
		SELECT production_id, TIMESTAMP(time_start) as time_start, TIMESTAMP(time_end) as time_end, cycles, is_completed 
		FROM rel_building_def_production
		WHERE building_id ?
		`, b.Id)
		err = rowp.Scan(&p.id, &p.timeStart, &p.timeEnd, &p.cycles, &p.isCompleted)
		if err != nil {
			http.Error(w, "could not find productions", 500)
			return
		}
		b.Production.Id = p.id
		b.Production.IsActive = !p.isCompleted
		b.Production.TimeEnd = p.timeEnd
		b.Production.Cycles.Total = p.cycles
		if p.isCompleted {
			b.Production.Cycles.Completed = p.cycles
		} else {
			// For duration get the duration from def_production
			pdQuery := db.QueryRow(`
			SELECT base_duration
			FROM def_production
			WHERE id ?
			`, p.id)
			var duration float64
			err = pdQuery.Scan(&duration)
			if err != nil {
				http.Error(w, "no production found", 500)
			}
			timePassed := time.Now().Unix() - p.timeStart.Unix()
			b.Production.Cycles.Completed = int(float64(timePassed) / duration)
		}

		//append buildingList
		buildingList = append(buildingList, b)
	}
	utils.SetHeaderJson(w)
	json.NewEncoder(w).Encode(buildingList)

}

func BuildingDetails(w http.ResponseWriter, req *http.Request) {
	utils.SetHeaderJson(w)
	id, err := utils.GetUrlParamId(utils.UrlParam{Url: req.URL.Path, Position: 3})
	if err != nil {
		http.Error(w, "could not find building id", http.StatusNotFound)
		return
	}
	claims, err := authentication.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	if !database.IsUserBuildingOwner(claims.UserId, id) {
		http.Error(w, "user is not building owner", http.StatusUnauthorized)
		return
	}
	//Get the building from the buildings Table
	building, err := database.FindBuilding(id)
	if err != nil {
		http.Error(w, "could not find building", http.StatusBadRequest)
		return
	}
	//Get the current Production of the Building
	production, err := database.FindProductionByBuilding(id)
	if err != nil {
		http.Error(w, "could not find productions", id)
	}

	storage, err := database.GetBStorage(id)

	type ResBody struct {
		database.Building
		Storage    []database.BStorage   `json:"storage"`
		Production []database.Production `json:"production"`
	}

	resBody := ResBody{
		Building:   building,
		Storage:    storage,
		Production: production,
	}
	json.NewEncoder(w).Encode(resBody)

}
