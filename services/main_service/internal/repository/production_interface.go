package repository

import "github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"

type ProductionRepository interface {
	GetProductionDefenitions() ([]model.ProductionDefenition, error)
	GetProductsOfProduction(productionID int) ([]model.ProductOfProduction, error)
	GetBuildingIdOfProduction(productionID int) (int, error)
	DeleteProductionByID(productionID int) error
	GetProductionsByBuildingType(buildingDefID int) ([]model.ProductionDefenition, error)
	AddProduction(newProduction model.BaseProduction) error
}
