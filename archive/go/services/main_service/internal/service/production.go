package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/repository"
)

type ProductionService interface {
	StartProduction(newProduction model.NewProductionRequest) error
	CancelProduction(productionID int) error
	IsProductionInBuilding(productionID int, buildingID int) error
	GetBuildingIdOfProduction(productionID int) (int, error)
}

type productionService struct {
	buildingRepo   repository.BuildingRepository
	producitonRepo repository.ProductionRepository
	userRepo       repository.UserRepository
}

func NewProductionService(buildingRepo repository.BuildingRepository, productionRepo repository.ProductionRepository, userRepo repository.UserRepository) ProductionService {
	return &productionService{
		buildingRepo:   buildingRepo,
		producitonRepo: productionRepo,
		userRepo:       userRepo,
	}
}

func (s *productionService) StartProduction(newPorduction model.NewProductionRequest) error {

	building, err := s.buildingRepo.GetBuildingByID(newPorduction.BuildingID)
	if err != nil {
		return fmt.Errorf("error finding building: %w", err)
	}

	bProductions, err := s.producitonRepo.GetProductionsByBuildingType(building.BuildingDefID)
	if err != nil {
		return fmt.Errorf("error finding productionDef: %w", err)
	}

	var prodDef model.ProductionDefenition
	for _, v := range bProductions {
		if newPorduction.ProductionID == v.ID {
			prodDef = v
		}

	}
	if prodDef.ID <= 0 {
		return fmt.Errorf("error no productionId match: %v", strconv.Itoa(prodDef.ID))
	}

	resources, err := s.userRepo.GetUserResources(building.UserID)

	if err != nil {
		return fmt.Errorf("no resources found: %w", err)
	}

	prodCost := prodDef.Cost * float64(newPorduction.Cycles)

	if prodCost > resources.Money {
		return fmt.Errorf("no sufficent ressources")
	}

	storage, err := s.buildingRepo.GetBuildingStorage(building.ID)

	if err != nil {
		return fmt.Errorf("erro finding storage: %w", err)
	}
	products, err := s.producitonRepo.GetProductsOfProduction(prodDef.ID)
	if err != nil {
		return fmt.Errorf("error finding products: %w", err)
	}

	for _, v := range products {
		v.Amount *= newPorduction.Cycles
		for _, k := range storage {
			if v.Amount <= k.Quantity {
				k.Quantity -= v.Amount
			} else {
				return fmt.Errorf("insufficent storage")
			}
		}
	}

	var newProduction model.BaseProduction = model.BaseProduction{
		ProductionID: prodDef.ID,
		BuildingID:   building.ID,
		Cycles:       newPorduction.Cycles,
		TimeStart:    time.Now().Unix(),
		TimeEnd:      time.Now().Unix() + int64(prodDef.BaseDuration+int64(newPorduction.Cycles)),
	}

	if err := s.producitonRepo.AddProduction(newProduction); err != nil {
		return fmt.Errorf("error starting production: %w", err)
	}

	return nil
}

func (s *productionService) CancelProduction(productionID int) error {
	return s.producitonRepo.DeleteProductionByID(productionID)
}

func (s *productionService) IsProductionInBuilding(productionID, buildingID int) error {
	id, err := s.producitonRepo.GetBuildingIdOfProduction(productionID)
	if err != nil {
		return fmt.Errorf("error finding buildingID: %w", err)
	}
	if id != buildingID {
		return fmt.Errorf("buildingIDs do not match: %w", err)
	}
	return nil
}

func (s *productionService) GetBuildingIdOfProduction(productionID int) (int, error) {
	return s.producitonRepo.GetBuildingIdOfProduction(productionID)
}
