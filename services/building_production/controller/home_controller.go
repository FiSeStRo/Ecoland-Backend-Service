package controller

import (
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

type HomeController struct {
	renderer *view.TemplateRenderer
}

func NewHomeController(renderer *view.TemplateRenderer) *HomeController {
	return &HomeController{
		renderer: renderer,
	}
}

func (c *HomeController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /config/storage", c.SaveToStorage)
}

func (c *HomeController) Index(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	data := map[string]any{
		"Title": "Building Production System",
	}

	c.renderer.Render(w, "home.html", data)
}

func (c *HomeController) SaveToStorage(w http.ResponseWriter, req *http.Request) {
	return
}
