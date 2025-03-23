package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
)

type ProductionRepository interface {
	GetDefProductions() ([]model.Production, error)
	CreateDefProduction(production model.Production) error
	GetProductionByProductID(productID int) ([]int, error)
}

type productionRepository struct {
	db *sql.DB
}

func NewProductionRepository(db *sql.DB) *productionRepository {
	return &productionRepository{
		db: db,
	}
}

func (r *productionRepository) GetDefProductions() ([]model.Production, error) {

	productionsQuery := `SELECT * FROM def_production`

	rows, err := r.db.Query(productionsQuery)
	if err != nil {
		return nil, fmt.Errorf("could not query def production rows: %w", err)
	}
	defer rows.Close()

	var result []model.Production

	for rows.Next() {
		var production model.Production
		if err := rows.Scan(&production.ID, &production.Name, &production.Cost, &production.Duration); err != nil {
			return nil, fmt.Errorf("error reading def_production: %w", err)
		}
		result = append(result, production)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	productQuery := `
	SELECT *
	FROM def_rel_production_product
	WHERE production_id IN (SELECT id FROM def_production)
	`

	productRow, err := r.db.Query(productQuery)
	if err != nil {
		return nil, fmt.Errorf("product relation query failed: %w", err)
	}
	defer productRow.Close()

	inputByProduction := make(map[int][]model.RelProduct)
	outputByProduction := make(map[int][]model.RelProduct)

	for productRow.Next() {
		var productionID int
		var isInput bool
		var relProduct model.RelProduct
		if err := productRow.Scan(&productionID, &relProduct.ProductID, &isInput, &relProduct.Amount); err != nil {
			return nil, fmt.Errorf("error scanning productRow: %w", err)
		}

		if isInput {
			inputByProduction[productionID] = append(inputByProduction[productionID], relProduct)
		} else {
			outputByProduction[productionID] = append(outputByProduction[productionID], relProduct)
		}
	}

	if err := productRow.Err(); err != nil {
		return nil, fmt.Errorf("product rows iteration failed: %w", err)
	}

	for i := range result {
		productionID := result[i].ID
		if input, exists := inputByProduction[productionID]; exists {
			result[i].InputType = input
		}
		if output, exists := outputByProduction[productionID]; exists {
			result[i].OutputType = output
		}
	}
	return result, nil
}

func (r *productionRepository) CreateDefProduction(production model.Production) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	defProducitonQuery := `INSERT INTO def_production(token_name, cost, base_duration) VALUES(?,?,?)`

	defProductionResult, err := tx.Exec(defProducitonQuery, production.Name, production.Cost, production.Duration)
	if err != nil {
		return fmt.Errorf("failed to insert def production: %w", err)
	}

	defProductionID, err := defProductionResult.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last inserted ID: %w", err)
	}

	relationStmt, err := tx.Prepare(`INSERT INTO def_rel_production_product(production_id, product_id, is_input, amount) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("could not prepare production product relation: %w", err)
	}
	defer relationStmt.Close()
	if len(production.InputType) > 0 {
		for _, input := range production.InputType {
			_, err := relationStmt.Exec(defProductionID, input.ProductID, true, input.Amount)
			if err != nil {
				return fmt.Errorf("failed to insert product production reltaion (production_id: %d, product_id: %d): %w",
					defProductionID, input.ProductID, err)
			}
		}
	}
	if len(production.OutputType) > 0 {
		for _, output := range production.OutputType {
			_, err := relationStmt.Exec(defProductionID, output.ProductID, false, output.Amount)
			if err != nil {
				return fmt.Errorf("failed to insert product production reltaion (production_id: %d, product_id: %d): %w",
					defProductionID, output.ProductID, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *productionRepository) GetProductionByProductID(productID int) ([]int, error) {
	query := `SELECT production_id FROM def_rel_production_product WHERE product_id = ?`
	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("could not get production IDs: %w", err)
	}
	var productionIDs []int
	for rows.Next() {
		var productionID int
		if err := rows.Scan(&productionID); err != nil {
			return nil, fmt.Errorf("could not read production IDs: %w", err)
		}
		productionIDs = append(productionIDs, productionID)
	}
	return productionIDs, nil
}
