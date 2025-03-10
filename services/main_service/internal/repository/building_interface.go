package repository

import (
	"database/sql"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
)

// BuildingRepository defines the interface for building-related data access
type BuildingRepository interface {
	GetBuildingDefinitions() ([]model.BuildingDefinition, error)
	GetBuildingDefinitionByID(id int) (*model.BuildingDefinition, error)
	GetBuildingsByUserID(userID int) ([]model.Building, error)
	GetBuildingByID(id int) (*model.Building, error)
	CreateBuilding(building *model.Building, tx *sql.Tx) (int, error)
	UpdateBuilding(building *model.Building) error
	GetBuildingStorage(buildingID int) ([]model.Storage, error)
	GetProductionByBuildingID(buildingID int) ([]model.Production, error)
	CheckBuildingOwnership(userID, buildingID int) (bool, error)
	BeginTx() (*sql.Tx, error)
}
