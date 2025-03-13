package service

import (
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/repository/mariadb"
)

type ProductionService interface {
	GetProductions() ([]model.Production, error)
	NewProduction(production model.Production) error
}

type productionService struct {
	productionRepository mariadb.ProductionRepository
}

func NewProductionService(productionRepository mariadb.ProductionRepository) ProductionService {
	return &productionService{
		productionRepository: productionRepository,
	}
}

func (s *productionService) GetProductions() ([]model.Production, error) {
	return s.productionRepository.GetDefProductions()
}

func (s *productionService) NewProduction(production model.Production) error {

	return nil
}
