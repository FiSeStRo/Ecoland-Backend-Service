package service

import (
	"fmt"
	"time"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/repository"
)

// BuildingService defines the interface for building-related operations
type BuildingService interface {
	GetConstructionList(userId int) ([]model.BuildingDefinition, error)
	ConstructBuilding(userId, buildingDefId int) (*model.ConstructResult, error)
	GetUserBuildings(userId int) ([]model.Building, error)
	GetBuildingDetails(userId, buildingId int) (*model.BuildingDetails, error)
	IsBuildingOwner(userID, buildingID int) (bool, error)
}

// buildingService implements BuildingService
type buildingService struct {
	buildingRepo repository.BuildingRepository
	userRepo     repository.UserRepository
}

// NewBuildingService creates a new building service
func NewBuildingService(buildingRepo repository.BuildingRepository, userRepo repository.UserRepository) BuildingService {
	return &buildingService{
		buildingRepo: buildingRepo,
		userRepo:     userRepo,
	}
}

// GetConstructionList returns buildings available for construction
func (s *buildingService) GetConstructionList(userId int) ([]model.BuildingDefinition, error) {
	// Get all building definitions
	definitions, err := s.buildingRepo.GetBuildingDefinitions()
	if err != nil {
		return nil, fmt.Errorf("failed to get building definitions: %w", err)
	}

	// Get user resources to determine what can be built
	resources, err := s.userRepo.GetUserResources(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user resources: %w", err)
	}

	// Enrich with availability based on user resources
	for i := range definitions {
		definitions[i].CanBuild = resources.Money >= definitions[i].BaseCost
	}

	return definitions, nil
}

// ConstructBuilding handles the business logic for constructing a new building
func (s *buildingService) ConstructBuilding(userId, buildingDefId int) (*model.ConstructResult, error) {
	// Get building definition
	buildingDef, err := s.buildingRepo.GetBuildingDefinitionByID(buildingDefId)
	if err != nil {
		return nil, fmt.Errorf("building definition not found: %w", err)
	}

	// Get user resources
	resources, err := s.userRepo.GetUserResources(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user resources: %w", err)
	}

	// Check if user has enough resources
	if resources.Money < buildingDef.BaseCost {
		return nil, fmt.Errorf("insufficient funds: %.2f required, %.2f available",
			buildingDef.BaseCost, resources.Money)
	}

	// Calculate completion time
	completionTime := time.Now().Add(time.Duration(buildingDef.BaseTime) * time.Minute)

	// Start transaction
	tx, err := s.buildingRepo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Deduct resources
	resources.Money -= buildingDef.BaseCost
	if err := s.userRepo.UpdateUserResources(resources); err != nil {
		return nil, fmt.Errorf("failed to update user resources: %w", err)
	}

	// Create building record
	building := &model.Building{
		UserID:                userId,
		BuildingDefID:         buildingDefId,
		Status:                "under_construction",
		ConstructionStartTime: time.Now(),
		ConstructionEndTime:   completionTime,
	}

	buildingID, err := s.buildingRepo.CreateBuilding(building, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create building: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &model.ConstructResult{
		BuildingID:     buildingID,
		CompletionTime: completionTime,
		RemainingFunds: resources.Money,
	}, nil
}

// GetUserBuildings returns all buildings owned by the user
func (s *buildingService) GetUserBuildings(userId int) ([]model.Building, error) {
	buildings, err := s.buildingRepo.GetBuildingsByUserID(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user buildings: %w", err)
	}

	// Check if any buildings have completed construction
	now := time.Now()
	for i, building := range buildings {
		if building.Status == "under_construction" && now.After(building.ConstructionEndTime) {
			// Update building status
			building.Status = "operational"
			err := s.buildingRepo.UpdateBuilding(&building)
			if err != nil {
				// Log error but continue
				fmt.Printf("failed to update building %d status: %v\n", building.ID, err)
			}
			buildings[i] = building
		}
	}

	return buildings, nil
}

// GetBuildingDetails returns details for a specific building
func (s *buildingService) GetBuildingDetails(userId, buildingId int) (*model.BuildingDetails, error) {
	// Get basic building info
	building, err := s.buildingRepo.GetBuildingByID(buildingId)
	if err != nil {
		return nil, fmt.Errorf("building not found: %w", err)
	}

	// Check ownership
	if building.UserID != userId {
		return nil, fmt.Errorf("unauthorized access to building")
	}

	// Get building definition
	buildingDef, err := s.buildingRepo.GetBuildingDefinitionByID(building.BuildingDefID)
	if err != nil {
		return nil, fmt.Errorf("building definition not found: %w", err)
	}

	// Get storage information
	storage, err := s.buildingRepo.GetBuildingStorage(buildingId)
	if err != nil {
		// Just log - this might be a new building with no storage
		fmt.Printf("failed to get building storage: %v\n", err)
	}

	// Get production information if building is operational
	var production []model.Production
	if building.Status == "operational" {
		production, err = s.buildingRepo.GetProductionByBuildingID(buildingId)
		if err != nil {
			// Just log - this might be a building with no production
			fmt.Printf("failed to get building production: %v\n", err)
		}
	}

	// Assemble complete building details
	details := &model.BuildingDetails{
		Building:   *building,
		Definition: *buildingDef,
		Storage:    storage,
		Production: production,
	}

	return details, nil
}

func (s *buildingService) IsBuildingOwner(userID, buildingID int) (bool, error) {
	isOwner, err := s.buildingRepo.CheckBuildingOwnership(userID, buildingID)
	if err != nil {
		return false, fmt.Errorf("failed to verify building ownership: %w", err)
	}

	return isOwner, nil
}
