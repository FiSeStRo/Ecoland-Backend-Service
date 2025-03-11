package controller

import (
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

// BuildingController handles building-related requests
type BuildingController struct {
	renderer *view.TemplateRenderer
}

// NewBuildingController creates a new building controller
func NewBuildingController(renderer *view.TemplateRenderer) *BuildingController {
	return &BuildingController{
		renderer: renderer,
	}
}

// Index handles the building page request
func (c *BuildingController) Index(w http.ResponseWriter, r *http.Request) {
	buildings := model.GetAllBuildings()

	data := map[string]interface{}{
		"Title":     "Buildings",
		"Buildings": buildings,
	}

	c.renderer.Render(w, "building.html", data)
}
