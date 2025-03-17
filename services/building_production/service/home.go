package service

import (
	"fmt"
	"path/filepath"

	"github.com/FiSeStRo/Ecoland-Backend-Service/go_pkg/config"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/repository/mariadb"
)

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
	storageDir := "../../file_storage"

	buildings, err := s.buildingRepo.GetDefBuildings()

	if err != nil {
		return fmt.Errorf("couldn't get buildings: %w", err)
	}
	productions, err := s.productionRepo.GetDefProductions()

	if err != nil {
		return fmt.Errorf("could not get productions: %w", err)
	}
	products, err := s.productRepo.GetDefProducts()

	if err != nil {
		return fmt.Errorf("could not get products: %w", err)
	}

	if err := config.WriteJSONFile(filepath.Join(storageDir, "def_buildings.json"), buildings); err != nil {
		return fmt.Errorf("failed to save buildings: %w", err)
	}

	if err := config.WriteJSONFile(filepath.Join(storageDir, "def_production.json"), productions); err != nil {
		return fmt.Errorf("failed to save productions: %w", err)
	}

	if err := config.WriteJSONFile(filepath.Join(storageDir, "def_products.json"), products); err != nil {
		return fmt.Errorf("failed to save products: %w", err)
	}
	return nil
}
