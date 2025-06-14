package service

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/FiSeStRo/Ecoland-Backend-Service/go_pkg/config"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/repository/mariadb"
)

type GetOverviewType struct {
	Buildings []struct {
		model.Building
		ProductIDs []int
	}
	Productions []struct {
		model.Production
		BuildingIDs []int
	}
	Products []struct {
		model.Product
		ProductionIDs []int
		BuildingIDs   []int
	}
}

type HomeService interface {
	SaveToStorage() error
	GetOverview() (GetOverviewType, error)
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

func (s *homeService) GetOverview() (GetOverviewType, error) {

	buildings, err := s.buildingRepo.GetDefBuildings()
	if err != nil {
		return GetOverviewType{}, fmt.Errorf("couldn't get buildings: %w", err)
	}
	productions, err := s.productionRepo.GetDefProductions()
	if err != nil {
		return GetOverviewType{}, fmt.Errorf("could not get productions: %w", err)
	}
	products, err := s.productRepo.GetDefProducts()
	if err != nil {
		return GetOverviewType{}, fmt.Errorf("could not get products: %w", err)
	}

	var productionRelations = make(map[int]struct {
		BuildingIDs []int
		ProductIDs  []int
	})
	for _, production := range productions {
		buildingIDs, err := s.buildingRepo.GetBuildingIDByProductionID(production.ID)
		if err != nil {
			log.Println(err)
			return GetOverviewType{}, fmt.Errorf("could not get buildignIds")
		}
		var productIDs []int
		for _, product := range production.InputType {
			productIDs = append(productIDs, product.ProductID)
		}

		for _, product := range production.OutputType {
			productIDs = append(productIDs, product.ProductID)
		}
		productionRelations[production.ID] = struct {
			BuildingIDs []int
			ProductIDs  []int
		}{
			BuildingIDs: buildingIDs,
			ProductIDs:  productIDs,
		}
	}
	var buildingsL []struct {
		model.Building
		ProductIDs []int
	}

	for _, building := range buildings {
		var productIDs []int
		for _, productionID := range building.Productions {
			if relation, exists := productionRelations[productionID]; exists {
				productIDs = append(productIDs, relation.ProductIDs...)
			}
		}
		buildingsL = append(buildingsL, struct {
			model.Building
			ProductIDs []int
		}{Building: building,
			ProductIDs: productIDs})

	}

	var productionsL []struct {
		model.Production
		BuildingIDs []int
	}

	for _, production := range productions {
		if relation, exists := productionRelations[production.ID]; exists {
			productionsL = append(productionsL, struct {
				model.Production
				BuildingIDs []int
			}{
				Production:  production,
				BuildingIDs: relation.BuildingIDs,
			})
		}
	}
	var productL []struct {
		model.Product
		ProductionIDs []int
		BuildingIDs   []int
	}
	for _, product := range products {

		productionIDs, err := s.productionRepo.GetProductionByProductID(product.ID)
		if err != nil {
			return GetOverviewType{}, fmt.Errorf("could not find productiond by product ID, %w", err)
		}

		var buildingIDs []int
		for _, productionID := range productionIDs {
			buildingIDs = append(buildingIDs, productionRelations[productionID].BuildingIDs...)
		}
		productL = append(productL, struct {
			model.Product
			ProductionIDs []int
			BuildingIDs   []int
		}{
			Product:       product,
			ProductionIDs: productionIDs,
			BuildingIDs:   buildingIDs,
		})
	}

	return (GetOverviewType{
		Buildings:   buildingsL,
		Productions: productionsL,
		Products:    productL,
	}), nil
}
