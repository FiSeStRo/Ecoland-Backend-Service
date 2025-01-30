package main

import (
	"database/sql"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/database"
	"github.com/FiSeStRo/Ecoland-Backend-Service/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/utils"
	_ "github.com/go-sql-driver/mysql"
)

var defSetup = flag.Bool("ds", false, "run setup")

func main() {
	flag.Parse()

	//Start Database
	dsn := "maria:maria123@tcp(localhost:3306)/ecoland"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panicln("panic can't open db", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Panicln("Error connecting to the database:", err)
	} else {
		log.Println("Succesfully connected to db")
	}
	database.SetDB(db)
	if *defSetup {
		log.Println("Setting up def tables")
		//setup defTable from config
		var defBuilding []database.DefBuilding
		var defProducts []database.DefProduct
		var defProduction []database.DefProduction
		err = utils.SetupDefTable("../config/def_buildings.json", defBuilding, "def_buildings", "(id, token_name, base_construction_cost, base_construction_time)")
		if err != nil {
			log.Fatalln(err)
		}
		err = utils.SetupDefTable("./config/def_products.json", defProducts, "def_product", "(id, base_value, token_name)")
		if err != nil {
			log.Fatalln(err)
		}
		err = utils.SetupDefTable("./config/def_production.json", defProduction, "def_production", "(id, token_name, cost, base_duration)")
		if err != nil {
			log.Fatalln(err)
		}
	}

	//Handle Routing

	//authentication
	http.HandleFunc("/authentication/refresh-token", service.NewRefreshToken)
	http.HandleFunc("/authentication/sign-in", service.SignIn)
	http.HandleFunc("/authentication/sign-up", service.SignUp)

	//buildings
	http.HandleFunc("/buildings/constructionlist", service.ConstructionList)
	http.HandleFunc("/buildings/construct", service.ConstructBuilding)
	http.HandleFunc("/buildings/list", service.ListOfBuildings)

	//production
	http.HandleFunc("/production/list", service.ListOfProductions)
	http.HandleFunc("/production/start", service.StartProduction)
	http.HandleFunc("/production/cancel", service.CancelProduction)

	// user
	http.HandleFunc("/user/resources", service.GetUserResources)

	http.HandleFunc("/", root)
	http.ListenAndServe(":8081", nil)
}

// Root handles the functionality of the root rout
func root(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hey ho Let's Go")
}
