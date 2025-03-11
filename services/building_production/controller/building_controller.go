package controller

import (
	"encoding/json"
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

// BuildingController handles building-related requests
type BuildingController struct {
	renderer        *view.TemplateRenderer
	buildingService service.BuildingService
}

// NewBuildingController creates a new building controller
func NewBuildingController(renderer *view.TemplateRenderer, buildingService service.BuildingService) *BuildingController {
	return &BuildingController{
		renderer:        renderer,
		buildingService: buildingService,
	}
}

func (c *BuildingController) RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/building", c.Index)
	mux.HandleFunc("POST /api/building", c.AddBuilding)
}

// Index handles the building page request
func (c *BuildingController) Index(w http.ResponseWriter, req *http.Request) {
	buildings, error := c.buildingService.GetAllBuildings()

	if error != nil {
		//TODO: add rendering of an error Page instead of the error message
		http.Error(w, "could not get Buidlings", 500)
		return
	}

	data := map[string]any{
		"Title":     "Buildings",
		"Buildings": buildings,
	}

	c.renderer.Render(w, "building.html", data)
}

func (c *BuildingController) AddBuilding(w http.ResponseWriter, req *http.Request) {

	var buildigng model.Building

	json.NewDecoder(req.Body).Decode(&buildigng)
	if err := c.buildingService.CreateBuilding(buildigng); err != nil {
		http.Error(w, "building could not be added", 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
