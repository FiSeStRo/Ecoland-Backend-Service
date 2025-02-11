package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/FiSeStRo/Ecoland-Backend-Service/database"
	"github.com/FiSeStRo/Ecoland-Backend-Service/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/utils"
	_ "github.com/go-sql-driver/mysql"
)

var defSetup = flag.Bool("ds", false, "run setup")

func main() {
	flag.Parse()
	//Start Database
	//for localhost "maria:maria123@tcp(localhost:3306)/ecoland"
	//Load Env
	err := utils.LoadEnv(".env")
	if err != nil {
		log.Println("error with env", err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PW"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
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
	database.InitDatabaseTables()
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

	mux := http.NewServeMux()
	//authentication
	mux.HandleFunc("/authentication/refresh-token", service.NewRefreshToken)
	mux.HandleFunc("/authentication/sign-in", service.SignIn)
	mux.HandleFunc("/authentication/sign-up", service.SignUp)

	//buildings
	mux.HandleFunc("/buildings/constructionlist", service.ConstructionList)
	mux.HandleFunc("/buildings/construct", service.ConstructBuilding)
	mux.HandleFunc("/buildings/list", service.ListOfBuildings)
	mux.HandleFunc("/buildings/details", service.BuildingDetails)

	//production
	mux.HandleFunc("/production/list", service.ListOfProductions)
	mux.HandleFunc("/production/start", service.StartProduction)
	mux.HandleFunc("/production/cancel", service.CancelProduction)

	//transportation
	mux.HandleFunc("/transportation/shipment", service.ShipItems)
	// user
	mux.HandleFunc("/user/resources", service.GetUserResources)
	mux.HandleFunc("GET /user/info", service.GetUserInfo)
	mux.HandleFunc("PATCH /user/info", service.UpdateUserInfo)
	mux.HandleFunc("/health", service.HealthCheck)
	mux.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":8081", enableCors(mux)))
}

// Root handles the functionality of the root rout
func root(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hey ho Let's Go")
}

func enableCors(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow any localhost origin
		origin := r.Header.Get("Origin")
		if strings.HasPrefix(origin, "http://localhost:") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
