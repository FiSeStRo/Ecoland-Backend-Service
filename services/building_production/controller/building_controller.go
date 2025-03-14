package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

// BuildingController handles building-related requests and manages
// the interaction between the HTTP layer and building services.
type BuildingController struct {
	renderer        *view.TemplateRenderer
	buildingService service.BuildingService
}

// NewBuildingController creates a new building controller with the specified
// template renderer and building service.
//
// Parameters:
//   - renderer: The HTML template renderer used to display building views
//   - buildingService: The service that handles building business logic
//
// Returns:
//   - A configured BuildingController instance ready to handle HTTP requests
func NewBuildingController(renderer *view.TemplateRenderer, buildingService service.BuildingService) *BuildingController {
	return &BuildingController{
		renderer:        renderer,
		buildingService: buildingService,
	}
}

// RegisterRoutes registers all building-related HTTP routes with the provided
// HTTP multiplexer (router).
//
// Parameters:
//   - mux: The HTTP multiplexer to register routes with
func (c *BuildingController) RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/building", c.Index)
	mux.HandleFunc("POST /api/buildings", c.AddBuilding)
}

// Index handles the building page request, displaying all available buildings.
// It fetches building data from the service and renders the building template.
//
// Parameters:
//   - w: The HTTP response writer
//   - req: The HTTP request
//
// HTTP Status Codes:
//   - 200 OK: Successfully returns the buildings page
//   - 500 Internal Server Error: If there's an error retrieving buildings

func (c *BuildingController) Index(w http.ResponseWriter, req *http.Request) {
	buildings, err := c.buildingService.GetAllBuildings()
	log.Println(buildings)
	if err != nil {
		log.Printf("could not get buildings: %v", err)
		//TODO: add rendering of an error Page instead of the error message
		http.Error(w, "could not get Buildings", 500)
		return
	}

	data := map[string]any{
		"Title":     "Buildings",
		"Buildings": buildings,
	}

	c.renderer.Render(w, "building.html", data)
}

// AddBuilding handles the API request to create a new building.
// It deserializes the JSON request body into a Building model
// and forwards it to the building service for creation.
//
// Parameters:
//   - w: The HTTP response writer
//   - req: The HTTP request containing building data in JSON format
//
// Expected Request Format:
//   - JSON body with building properties (name, resourceCost, buildTime, productions)
//
// HTTP Status Codes:
//   - 201 Created: Successfully created the building
//   - 500 Internal Server Error: If there's an error creating the building

func (c *BuildingController) AddBuilding(w http.ResponseWriter, req *http.Request) {

	var buildigng model.Building

	json.NewDecoder(req.Body).Decode(&buildigng)
	log.Println("building", buildigng)
	if err := c.buildingService.CreateBuilding(buildigng); err != nil {
		http.Error(w, "building could not be added", 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
