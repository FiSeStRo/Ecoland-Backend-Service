package service

import "github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/repository/mariadb"

type HomeService interface {
	SaveToStorage() error
}

type homeService struct {
	buildingRepo   mariadb.BuildingRepository
	productionRepo mariadb.ProductionRepository
	productRepo    mariadb.ProductRepository
}

func NewHomeService(buildingRepo mariadb.BuildingRepository, productionRepo mariadb.ProductionRepository, productRepo mariadb.ProductRepository) HomeService {
	return &homeService{
		buildingRepo:   buildingRepo,
		productionRepo: productionRepo,
		productRepo:    productRepo,
	}
}

func (s *homeService) SaveToStorage() error {

	//TODO: get Data from repos

	//TODO: save Data to json
	return nil
}
