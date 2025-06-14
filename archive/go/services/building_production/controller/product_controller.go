package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

// ProductController handles product-related requests
type ProductController struct {
	renderer       *view.TemplateRenderer
	productService service.ProductService
}

// NewProductController creates a new product controller
func NewProductController(renderer *view.TemplateRenderer, productService service.ProductService) *ProductController {
	return &ProductController{
		renderer:       renderer,
		productService: productService,
	}
}

// Index handles the product page request
func (c *ProductController) Index(w http.ResponseWriter, r *http.Request) {
	products, err := c.productService.GetDefProducts()

	if err != nil {
		log.Println("could not get products")
		http.Error(w, "error getting products", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":    "Products",
		"Products": products,
	}

	c.renderer.Render(w, "product.html", data)
}

// RegiserRoutes handles the routes registered by the ProductController
func (c *ProductController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/product", c.Index)
	mux.HandleFunc("POST api/product", c.NewProduct)
}

// NewProduct controls the request of creating a new product
func (c *ProductController) NewProduct(w http.ResponseWriter, req *http.Request) {

	var product model.Product

	if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
		log.Println("error decoding new product request: %w", err)
		http.Error(w, "error with request body", http.StatusBadRequest)
		return
	}

	if err := c.productService.AddProduct(product); err != nil {
		http.Error(w, "error adding product to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
