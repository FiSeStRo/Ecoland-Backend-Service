package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/FiSeStRo/Ecoland-Backend-Service/go_pkg/database"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/controller"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/repository/mariadb"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

var magrationFlag bool

func init() {
	flag.BoolVar(&magrationFlag, "migrate", false, "Run config migration")
}

func main() {

	flag.Parse()

	renderer, err := view.NewTemplateRenderer("view/templates")
	if err != nil {
		log.Fatalf("Failed to initialize template renderer: %v", err)
	}

	dbConfig := database.NewConfig()

	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	if magrationFlag {
		if err = database.MigrateDefinitonData(db); err != nil {
			log.Println("could not migrate Data:", err)
		}
	}

	buildingRepo := mariadb.NewBuildingRepository(db)
	productionRepo := mariadb.NewProductionRepository(db)
	productRepo := mariadb.NewProductRepository(db)

	buildingService := service.NewBuildingService(buildingRepo)
	productionService := service.NewProductionService(productionRepo)
	productService := service.NewProductService(productRepo)
	homeService := service.NewHomeService(buildingRepo, productionRepo, productRepo)

	homeController := controller.NewHomeController(renderer, homeService)
	buildingController := controller.NewBuildingController(renderer, buildingService)
	productionController := controller.NewProductionController(renderer, productionService)
	productController := controller.NewProductController(renderer, productService)

	mux := http.NewServeMux()

	// Register routes
	homeController.RegisterRoutes(mux)
	buildingController.RegisterRoutes(mux)
	productionController.RegisterRoutes(mux)
	productController.RegisterRoutes(mux)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	const port = 8082

	log.Println("Starting building and production server on:", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), mux))
}
