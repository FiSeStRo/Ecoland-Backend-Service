package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/repository"
)

// buildingRepository implements the BuildingRepository interface
type buildingRepository struct {
	db *sql.DB
}

// NewBuildingRepository creates a new building repository
func NewBuildingRepository(db *sql.DB) repository.BuildingRepository {
	return &buildingRepository{db: db}
}

// GetBuildingDefinitions returns all available building definitions
func (r *buildingRepository) GetBuildingDefinitions() ([]model.BuildingDefinition, error) {
	query := `SELECT id, token_name, base_construction_cost, base_construction_time 
              FROM def_buildings`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var definitions []model.BuildingDefinition
	for rows.Next() {
		var def model.BuildingDefinition
		if err := rows.Scan(&def.ID, &def.Name, &def.BaseCost, &def.BaseTime); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		definitions = append(definitions, def)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return definitions, nil
}

// GetBuildingDefinitionByID returns a building definition by ID
func (r *buildingRepository) GetBuildingDefinitionByID(id int) (*model.BuildingDefinition, error) {
	query := `SELECT id, token_name, base_construction_cost, base_construction_time 
              FROM def_buildings 
              WHERE id = ?`

	var def model.BuildingDefinition
	err := r.db.QueryRow(query, id).Scan(&def.ID, &def.Name, &def.BaseCost, &def.BaseTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("building definition not found with ID %d", id)
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &def, nil
}

// GetBuildingsByUserID returns all buildings owned by a user
func (r *buildingRepository) GetBuildingsByUserID(userID int) ([]model.Building, error) {
	query := `SELECT id, user_id, building_def_id, status, construction_start_time, construction_end_time 
              FROM buildings 
              WHERE user_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var buildings []model.Building
	for rows.Next() {
		var b model.Building
		var startTime, endTime sql.NullTime

		if err := rows.Scan(&b.ID, &b.UserID, &b.BuildingDefID, &b.Status, &startTime, &endTime); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		if startTime.Valid {
			b.ConstructionStartTime = startTime.Time
		}

		if endTime.Valid {
			b.ConstructionEndTime = endTime.Time
		}

		buildings = append(buildings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return buildings, nil
}

// GetBuildingByID returns a specific building by ID
func (r *buildingRepository) GetBuildingByID(id int) (*model.Building, error) {
	query := `SELECT id, user_id, building_def_id, status, construction_start_time, construction_end_time 
              FROM buildings 
              WHERE id = ?`

	var b model.Building
	var startTime, endTime sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&b.ID, &b.UserID, &b.BuildingDefID, &b.Status, &startTime, &endTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("building not found with ID %d", id)
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	if startTime.Valid {
		b.ConstructionStartTime = startTime.Time
	}

	if endTime.Valid {
		b.ConstructionEndTime = endTime.Time
	}

	return &b, nil
}

// CreateBuilding creates a new building record
func (r *buildingRepository) CreateBuilding(building *model.Building, tx *sql.Tx) (int, error) {
	query := `INSERT INTO buildings 
              (user_id, building_def_id, status, construction_start_time, construction_end_time) 
              VALUES (?, ?, ?, ?, ?)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.Exec(
			query,
			building.UserID,
			building.BuildingDefID,
			building.Status,
			building.ConstructionStartTime,
			building.ConstructionEndTime,
		)
	} else {
		result, err = r.db.Exec(
			query,
			building.UserID,
			building.BuildingDefID,
			building.Status,
			building.ConstructionStartTime,
			building.ConstructionEndTime,
		)
	}

	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return int(id), nil
}

// UpdateBuilding updates an existing building
func (r *buildingRepository) UpdateBuilding(building *model.Building) error {
	query := `UPDATE buildings 
              SET status = ?, construction_start_time = ?, construction_end_time = ? 
              WHERE id = ?`

	_, err := r.db.Exec(
		query,
		building.Status,
		building.ConstructionStartTime,
		building.ConstructionEndTime,
		building.ID,
	)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}

// GetBuildingStorage returns the storage contents of a building
func (r *buildingRepository) GetBuildingStorage(buildingID int) ([]model.Storage, error) {
	query := `SELECT s.id, s.building_id, s.product_def_id, s.quantity, p.token_name 
              FROM storage s
              JOIN def_products p ON s.product_def_id = p.id
              WHERE s.building_id = ?`

	rows, err := r.db.Query(query, buildingID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var storageItems []model.Storage
	for rows.Next() {
		var item model.Storage
		if err := rows.Scan(&item.ID, &item.BuildingID, &item.ProductDefID, &item.Quantity, &item.ProductName); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		storageItems = append(storageItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return storageItems, nil
}

// GetProductionByBuildingID returns the production configurations for a building
func (r *buildingRepository) GetProductionByBuildingID(buildingID int) ([]model.Production, error) {
	query := `SELECT p.id, p.building_id, p.product_def_id, p.status, p.production_rate, dp.token_name
              FROM production p
              JOIN def_products dp ON p.product_def_id = dp.id
              WHERE p.building_id = ?`

	rows, err := r.db.Query(query, buildingID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var productions []model.Production
	for rows.Next() {
		var prod model.Production
		if err := rows.Scan(&prod.ID, &prod.BuildingID, &prod.ProductDefID, &prod.Status, &prod.Rate, &prod.ProductName); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		productions = append(productions, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return productions, nil
}

// CheckBuildingOwnership verifies if a user owns a specific building
func (r *buildingRepository) CheckBuildingOwnership(userID, buildingID int) (bool, error) {
	// Get the building
	building, err := r.GetBuildingByID(buildingID)

	// Handle not found error
	if err != nil {
		if err.Error() == fmt.Sprintf("building not found with ID %d", buildingID) {
			return false, nil // Building doesn't exist
		}
		return false, fmt.Errorf("error checking building ownership: %w", err)
	}

	// Check if the user ID matches the building's owner ID
	return building.UserID == userID, nil
}

// BeginTx starts a new transaction
func (r *buildingRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}
