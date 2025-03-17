package controller

import (
	"log"
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

type HomeController struct {
	renderer    *view.TemplateRenderer
	homeService service.HomeService
}

func NewHomeController(renderer *view.TemplateRenderer, homeService service.HomeService) *HomeController {
	return &HomeController{
		renderer:    renderer,
		homeService: homeService,
	}
}

func (c *HomeController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", c.Index)
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

	if err := c.homeService.SaveToStorage(); err != nil {
		log.Println("error saving files: ", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
