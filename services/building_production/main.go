package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/controller"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

func main() {

	renderer, err := view.NewTemplateRenderer("view/templates")
	if err != nil {
		log.Fatalf("Failed to initialize template renderer: %v", err)
	}

	homeController := controller.NewHomeController(renderer)
	buildingController := controller.NewBuildingController(renderer)
	productController := controller.NewProductController(renderer)
	productionController := controller.NewProductionController(renderer)

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", homeController.Index)
	mux.HandleFunc("/building", buildingController.Index)
	mux.HandleFunc("/production", productionController.Index)
	mux.HandleFunc("/product", productController.Index)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	const port = 8082

	log.Println("Starting buildigng and production server on:", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), mux))
}
