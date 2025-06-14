package service

import (
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/repository/mariadb"
)

type BuildingService interface {
	GetAllBuildings() ([]model.Building, error)
	CreateBuilding(building model.Building) error
}

type buildingService struct {
	buildingRepo mariadb.BuildingRepository
}

func NewBuildingService(buildingRepo mariadb.BuildingRepository) BuildingService {
	return &buildingService{
		buildingRepo: buildingRepo,
	}
}

func (s *buildingService) GetAllBuildings() ([]model.Building, error) {
	return s.buildingRepo.GetDefBuildings()
}

func (s *buildingService) CreateBuilding(building model.Building) error {
	return s.buildingRepo.CreateDefBuilding(building)
}
