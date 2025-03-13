package controller

import (
	"net/http"
	"strconv"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/view"
)

// ProductionController handles production-related requests
type ProductionController struct {
	renderer          *view.TemplateRenderer
	productionService service.ProductionService
}

// NewProductionController creates a new production controller
func NewProductionController(renderer *view.TemplateRenderer, productionService service.ProductionService) *ProductionController {
	return &ProductionController{
		renderer:          renderer,
		productionService: productionService,
	}
}

// Index handles the production page request
func (c *ProductionController) Index(w http.ResponseWriter, req *http.Request) {
	productions, err := c.productionService.GetProductions()

	if err != nil {
		http.Error(w, "failed to get productions", 500)
		return
	}
	type ProductionWText struct {
		model.Production
		OutputText string
		InputText  string
	}

	var productionsWText []ProductionWText
	for _, prod := range productions {
		var inputText string
		var outputText string
		for i, input := range prod.InputType {
			text := strconv.Itoa(input.ProductID) + ":" + strconv.Itoa(input.Amount)
			if i == 0 {
				inputText = text
			} else {
				inputText += "," + text
			}
		}
		for i, output := range prod.OutputType {
			text := strconv.Itoa(output.ProductID) + ":" + strconv.Itoa(output.Amount)
			if i == 0 {
				outputText = text
			} else {
				outputText += "," + text
			}
		}
		productionsWText = append(productionsWText, ProductionWText{Production: prod, InputText: inputText, OutputText: outputText})
	}
	data := map[string]any{
		"Title":       "Production",
		"Productions": productionsWText,
	}

	c.renderer.Render(w, "production.html", data)
}

func (c *ProductionController) AddProduction(w http.ResponseWriter, req *http.Request) {

}
