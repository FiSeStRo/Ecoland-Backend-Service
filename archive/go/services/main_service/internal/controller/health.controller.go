package controller

import (
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/utils"
)

type HealthController struct {
	healthService service.HealthService
}

func NewHealthController(healthService service.HealthService) *HealthController {
	return &HealthController{
		healthService: healthService,
	}
}

func (c *HealthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", c.HealthStatus)
}

func (c *HealthController) HealthStatus(w http.ResponseWriter, req *http.Request) {
	healthStatus := c.healthService.GetHealthStatus()
	type HealthStatus struct {
		Status string `json:"status"`
	}
	status := HealthStatus{Status: healthStatus}
	utils.WriteJSON(w, status, 200)
}
