package controller

import (
	"encoding/json"
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/utils"
)

type ProductionController struct {
	productionService service.ProductionService
	buildingService   service.BuildingService
	authService       service.AuthService
}

func NewProductionController(productionService service.ProductionService, authService service.AuthService, buildingService service.BuildingService) *ProductionController {
	return &ProductionController{
		productionService: productionService,
		buildingService:   buildingService,
		authService:       authService,
	}
}

func (c *ProductionController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /production/start", c.ProductionStart)
	mux.HandleFunc("DELETE /production/cancel", c.ProductionCancel)
}

func (c *ProductionController) ProductionStart(w http.ResponseWriter, req *http.Request) {

	claims, err := c.authService.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	var reqBody model.NewProductionRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "malformed request", 400)
		return
	}

	if ok, err := c.buildingService.IsBuildingOwner(claims.UserId, reqBody.BuildingID); !ok || err != nil {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}
	if err := c.productionService.StartProduction(reqBody); err != nil {
		http.Error(w, "Internal server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (c *ProductionController) ProductionCancel(w http.ResponseWriter, req *http.Request) {
	claims, err := c.authService.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "AccessDenied", http.StatusForbidden)
		return
	}

	productionID, err := utils.GetUrlParamId(utils.UrlParam{Url: req.URL.Path, Position: 3})
	if err != nil {
		http.Error(w, "id Param not found", 400)
		return
	}
	buildingID, err := c.productionService.GetBuildingIdOfProduction(productionID)
	if err != nil {
		http.Error(w, "Sorry something went wrong", 500)
		return
	}

	if ok, err := c.buildingService.IsBuildingOwner(claims.UserId, buildingID); !ok || err != nil {
		http.Error(w, "AccessDenied", http.StatusForbidden)
		return
	}

	if err := c.productionService.CancelProduction(productionID); err != nil {
		http.Error(w, "could not cancel Production", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
