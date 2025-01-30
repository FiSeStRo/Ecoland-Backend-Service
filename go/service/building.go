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

	stmt, err := db.Prepare(`INSERT INTO buildings (user_id, def_id, name, time_build)
	VALUES (?, ?, ?, CURRENT_TIMESTAMP)`)
	if err != nil {
		log.Println("Error preparing data insert", err)
		http.Error(w, "Error Preparing data insert", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(claims.UserId, constructReq.Def_id, constructReq.Display_name)
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
		//TODO: collect the rest of the data from the other tables and return the complete struct if success

		//get the list of possible productions
		// stmtps, err := db.Query(`
		// SELECT p.production_id, dp.token_name, dp.cost, dp.base_duration
		// FROM def_rel_building_production p
		// JOIN def_production dp ON p.production_id = dp.id
		// WHERE p.building_id=?`, b.id)
		// if err != nil {
		// 	http.Error(w, "error searching for productions", http.StatusInternalServerError)
		// 	return
		// }
		// for stmtps.Next() {
		// 	p := production{}
		// 	err = stmtps.Scan(&p.id, &p.token_name, &p.cost, &p.base_duration)
		// 	if err != nil {
		// 		http.Error(w, "error scan for productions", http.StatusInternalServerError)
		// 		return
		// 	}
		// 	b.productions = append(b.productions, p)
		// }
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
