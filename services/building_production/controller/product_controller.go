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
		log.Println("could not ger products")
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

func (c *ProductController) NewProduct(w http.ResponseWriter, req *http.Request) {

	var product model.Product

	if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
		log.Println("error decoding new product request: %w", err)
		http.Error(w, "error with request body", http.StatusBadRequest)
		return
	}

	//TODO: decode data
	//TODO: add service
	//TODO: Return response if created
}

//TODO: add product repository to getProducts and add a Product
//TODO: add the addProduct Controller funct and registerRoutes
