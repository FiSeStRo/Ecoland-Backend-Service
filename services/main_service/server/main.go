package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/FiSeStRo/Ecoland-Backend-Service/go/pkg/database"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/config"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/controller"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/repository/mariadb"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/utils"
)

var migrateFlag bool

func init() {
	flag.BoolVar(&migrateFlag, "migrate", false, "Run database migrations")
}

func main() {

	flag.Parse()

	if err := config.LoadEnv("./.env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// setup configurations
	jwtConfig := config.LoadJWTConfig()
	dbConfig := database.NewConfig()

	// Connect to database
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	if migrateFlag {
		log.Println("Running Database migrations...")
		if err := database.RunMigrations(db); err != nil {
			log.Panicln("Database migration Failed with: %w", err)
		}
		log.Println("Migration completed successfully")
	}

	// Initialize repositories
	userRepo := mariadb.NewUserRepository(db)
	buildingRepo := mariadb.NewBuildingRepository(db)
	productionRepo := mariadb.NewProductionRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	buildingService := service.NewBuildingService(buildingRepo, userRepo)
	authService := service.NewAuthService(jwtConfig, userRepo)
	productionService := service.NewProductionService(buildingRepo, productionRepo, userRepo)
	healthService := service.NewHealthService()

	// Initialize controllers
	userController := controller.NewUserController(userService, authService)
	buildingController := controller.NewBuildingController(buildingService, authService)
	productionController := controller.NewProductionController(productionService, authService, buildingService)
	healthController := controller.NewHealthController(healthService)

	// Setup router
	mux := http.NewServeMux()

	// Register routes
	userController.RegisterRoutes(mux)
	buildingController.RegisterRoutes(mux)
	productionController.RegisterRoutes(mux)
	healthController.RegisterRoutes(mux)

	// Start server
	const port = 8081
	log.Println("Starting server on :", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), utils.EnableCors(mux)))
}
