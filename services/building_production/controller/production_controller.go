package controller

import (
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

// ProductionController handles production-related requests
type ProductionController struct {
	renderer *view.TemplateRenderer
}

// NewProductionController creates a new production controller
func NewProductionController(renderer *view.TemplateRenderer) *ProductionController {
	return &ProductionController{
		renderer: renderer,
	}
}

// Index handles the production page request
func (c *ProductionController) Index(w http.ResponseWriter, r *http.Request) {
	productions := model.GetAllProductions()

	data := map[string]any{
		"Title":       "Production",
		"Productions": productions,
	}

	c.renderer.Render(w, "production.html", data)
}
