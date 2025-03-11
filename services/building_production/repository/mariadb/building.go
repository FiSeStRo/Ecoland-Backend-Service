package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
)

type BuildingRepository interface {
	GetDefBuildings() ([]model.Building, error)
	CreateDefBuilding(building model.Building) error
}

type buildingRepository struct {
	db *sql.DB
}

func NewBuildingRepository(db *sql.DB) *buildingRepository {
	return &buildingRepository{
		db: db,
	}
}

func (r *buildingRepository) GetDefBuildings() ([]model.Building, error) {
	buildingsQuery := `
        SELECT id, token_name, base_construction_cost, base_construction_time 
        FROM def_buildings`

	rows, err := r.db.Query(buildingsQuery)
	if err != nil {
		return nil, fmt.Errorf("buildings query failed: %w", err)
	}
	defer rows.Close()

	buildingsMap := make(map[int]*model.Building)
	var result []model.Building

	for rows.Next() {
		var building model.Building
		if err := rows.Scan(&building.ID, &building.Name, &building.ResourceCost, &building.BuildTime); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		building.Production = make([]int, 0)
		result = append(result, building)
		buildingsMap[building.ID] = &result[len(result)-1]
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	productionsQuery := `
        SELECT building_id, production_id 
        FROM rel_building_def_production 
        WHERE building_id IN (SELECT id FROM def_buildings)`

	prodRows, err := r.db.Query(productionsQuery)
	if err != nil {
		return nil, fmt.Errorf("production relations query failed: %w", err)
	}
	defer prodRows.Close()

	for prodRows.Next() {
		var buildingID, productionID int
		if err := prodRows.Scan(&buildingID, &productionID); err != nil {
			return nil, fmt.Errorf("production scan failed: %w", err)
		}
		if building, exists := buildingsMap[buildingID]; exists {
			building.Production = append(building.Production, productionID)
		}
	}

	if err := prodRows.Err(); err != nil {
		return nil, fmt.Errorf("production rows iteration failed: %w", err)
	}

	return result, nil
}

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
	buildingQuery := `INSERT INTO def_buildings(token_name, base_construction_cost, base_construction_time) VALUES (?, ?, ?)`
	buildingResult, err := tx.Exec(buildingQuery, building.Name, building.ResourceCost, building.BuildTime)
	if err != nil {
		return fmt.Errorf("failed to insert building: %w", err)
	}
	buildingID, err := buildingResult.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	if len(building.Production) > 0 {
		relationStmt, err := tx.Prepare(`INSERT INTO rel_building_def_production(building_id, production_id) VALUES (?, ?)`)
		if err != nil {
			return fmt.Errorf("failed to prepare relation statement: %w", err)
		}
		defer relationStmt.Close()

		for _, productionID := range building.Production {
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
