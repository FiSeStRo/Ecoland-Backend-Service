package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/FiSeStRo/Ecoland-Backend-Service/go_pkg/database"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/controller"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/repository/mariadb"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

func main() {

	renderer, err := view.NewTemplateRenderer("view/templates")
	if err != nil {
		log.Fatalf("Failed to initialize template renderer: %v", err)
	}

	dbConfig := database.NewConfig()

	// Connect to database
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	if err = database.MigrateDefinitonData(db); err != nil {
		log.Println("could not migrate Data:", err)
	}
	buildingRepo := mariadb.NewBuildingRepository(db)
	productionRepo := mariadb.NewProductionRepository(db)

	buildingService := service.NewBuildingService(buildingRepo)
	productionService := service.NewProductionService(productionRepo)

	homeController := controller.NewHomeController(renderer)
	buildingController := controller.NewBuildingController(renderer, buildingService)
	productionController := controller.NewProductionController(renderer, productionService)
	productController := controller.NewProductController(renderer)

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", homeController.Index)
	buildingController.RegisterRoutes(mux)
	productionController.RegisterRoutes(mux)
	mux.HandleFunc("/product", productController.Index)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	const port = 8082

	log.Println("Starting building and production server on:", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), mux))
}
