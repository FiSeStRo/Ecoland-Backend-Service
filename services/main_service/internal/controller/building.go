package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/service"
)

// BuildingController handles HTTP requests related to buildings
type BuildingController struct {
	buildingService service.BuildingService
	authService     service.AuthService
}

// NewBuildingController creates a new building controller
func NewBuildingController(buildingService service.BuildingService, authService service.AuthService) *BuildingController {
	return &BuildingController{
		buildingService: buildingService,
		authService:     authService,
	}
}

// RegisterRoutes registers all building-related routes
func (c *BuildingController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /buildings/constructionlist", c.ConstructionList)
	mux.HandleFunc("POST /buildings/construct", c.ConstructBuilding)
	mux.HandleFunc("GET /buildings/list", c.ListOfBuildings)
	mux.HandleFunc("GET /buildings/details", c.BuildingDetails)
}

// ConstructionList returns the list of buildings available for construction
func (c *BuildingController) ConstructionList(w http.ResponseWriter, req *http.Request) {
	claims, err := c.authService.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	buildingList, err := c.buildingService.GetConstructionList(claims.UserId)
	if err != nil {
		http.Error(w, "Failed to get construction list: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildingList)
}

// ConstructBuilding handles construction of a new building
func (c *BuildingController) ConstructBuilding(w http.ResponseWriter, req *http.Request) {

	claims, err := c.authService.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var constructionReq model.ConstructRequest
	if err := json.NewDecoder(req.Body).Decode(&constructionReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := c.buildingService.ConstructBuilding(claims.UserId, constructionReq.BuildingDefID)
	if err != nil {
		http.Error(w, "Failed to construct building: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ListOfBuildings returns all buildings owned by the user
func (c *BuildingController) ListOfBuildings(w http.ResponseWriter, req *http.Request) {
	claims, err := c.authService.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	buildings, err := c.buildingService.GetUserBuildings(claims.UserId)
	if err != nil {
		http.Error(w, "Failed to fetch buildings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildings)
}

// BuildingDetails returns details for a specific building
func (c *BuildingController) BuildingDetails(w http.ResponseWriter, req *http.Request) {
	buildingIdStr := req.URL.Query().Get("id")
	buildingId, err := strconv.Atoi(buildingIdStr)
	if err != nil || buildingId <= 0 {
		http.Error(w, "Invalid building ID", http.StatusBadRequest)
		return
	}

	claims, err := c.authService.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	buildingDetails, err := c.buildingService.GetBuildingDetails(claims.UserId, buildingId)
	if err != nil {
		http.Error(w, "Failed to get building details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildingDetails)
}
