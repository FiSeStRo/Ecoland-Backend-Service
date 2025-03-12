package controller

import (
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

// ProductController handles product-related requests
type ProductController struct {
	renderer *view.TemplateRenderer
}

// NewProductController creates a new product controller
func NewProductController(renderer *view.TemplateRenderer) *ProductController {
	return &ProductController{
		renderer: renderer,
	}
}

// Index handles the product page request
func (c *ProductController) Index(w http.ResponseWriter, r *http.Request) {
	products := model.GetAllProducts()

	data := map[string]any{
		"Title":    "Products",
		"Products": products,
	}

	c.renderer.Render(w, "product.html", data)
}
