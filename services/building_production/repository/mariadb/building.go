package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
)

// BuildingRepository defines the interface for building data storage operations.
// It provides methods for retrieving and creating building definitions.
type BuildingRepository interface {
	GetDefBuildings() ([]model.Building, error)
	CreateDefBuilding(building model.Building) error
	GetBuildingIDByProductionID(ID int) ([]int, error)
}

// buildingRepository implements the BuildingRepository interface
type buildingRepository struct {
	db *sql.DB
}

// NewBuildingRepository creates a new building repository with the provided database connection.
//
// Parameters:
//   - db: A valid SQL database connection to MariaDB
//
// Returns:
//   - A pointer to a buildingRepository that implements the BuildingRepository interface
func NewBuildingRepository(db *sql.DB) *buildingRepository {
	return &buildingRepository{
		db: db,
	}
}

// GetDefBuildings retrieves all building definitions from the database.
// It first queries basic building information, then retrieves and associates
// production relationships in a second query for optimized database access.
//
// Returns:
//   - A slice of Building models with their associated production IDs
//   - An error if any database operations fail
func (r *buildingRepository) GetDefBuildings() ([]model.Building, error) {
	buildingsQuery := `
        SELECT id, token_name, base_construction_cost, base_construction_time 
        FROM def_buildings`

	rows, err := r.db.Query(buildingsQuery)
	if err != nil {
		return nil, fmt.Errorf("buildings query failed: %w", err)
	}
	defer rows.Close()

	var result []model.Building

	for rows.Next() {
		var building model.Building
		if err := rows.Scan(&building.ID, &building.Name, &building.ResourceCost, &building.BuildTime); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		result = append(result, building)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}
	productionsQuery := `
	SELECT building_id, production_id 
	FROM def_rel_building_production
	WHERE building_id IN (SELECT id FROM def_buildings)`

	prodRows, err := r.db.Query(productionsQuery)
	if err != nil {
		return nil, fmt.Errorf("production relations query failed: %w", err)
	}
	defer prodRows.Close()

	productionsByBuilding := make(map[int][]int)

	for prodRows.Next() {
		var buildingID, productionID int
		if err := prodRows.Scan(&buildingID, &productionID); err != nil {
			return nil, fmt.Errorf("production scan failed: %w", err)
		}

		productionsByBuilding[buildingID] = append(productionsByBuilding[buildingID], productionID)
	}

	if err := prodRows.Err(); err != nil {
		return nil, fmt.Errorf("production rows iteration failed: %w", err)
	}

	for i := range result {
		buildingID := result[i].ID
		if prods, exists := productionsByBuilding[buildingID]; exists {
			result[i].Productions = prods
		}
	}

	return result, nil
}

// CreateDefBuilding creates a new building definition in the database.
// It uses a transaction to ensure that both the building and its production
// relationships are created atomically.
//
// Parameters:
//   - building: The Building model containing the building definition data and production IDs
//
// Returns:
//   - An error if the database operation fails at any point
//
// The function will automatically roll back the transaction if any step fails.
func (r *buildingRepository) CreateDefBuilding(building model.Building) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	buildingQuery := `INSERT INTO def_buildings( token_name, base_construction_cost, base_construction_time) VALUES ( ?, ?, ?)`
	buildingResult, err := tx.Exec(buildingQuery, building.Name, building.ResourceCost, building.BuildTime)
	if err != nil {
		return fmt.Errorf("failed to insert building: %w", err)
	}
	buildingID, err := buildingResult.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	if len(building.Productions) > 0 {
		relationStmt, err := tx.Prepare(`INSERT INTO def_rel_building_production(building_id, production_id) VALUES (?, ?)`)
		if err != nil {
			return fmt.Errorf("failed to prepare relation statement: %w", err)
		}
		defer relationStmt.Close()

		for _, productionID := range building.Productions {
			_, err = relationStmt.Exec(buildingID, productionID)
			if err != nil {
				return fmt.Errorf("failed to insert production relation (building_id=%d, production_id=%d): %w",
					buildingID, productionID, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *buildingRepository) GetBuildingIDByProductionID(ID int) ([]int, error) {
	query := `SELECT building_id FROM def_rel_building_production WHERE production_id = ?`

	rows, err := r.db.Query(query, ID)
	if err != nil {
		return nil, fmt.Errorf("could not get buildings by production id : %w", err)
	}
	var buidlingsID []int
	for rows.Next() {
		var buildingID int
		if err := rows.Scan(&buildingID); err != nil {
			return nil, fmt.Errorf("could not save buildings by production id : %w", err)
		}
		buidlingsID = append(buidlingsID, buildingID)
	}
	return buidlingsID, nil
}
