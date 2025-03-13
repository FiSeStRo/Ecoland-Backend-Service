package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
)

type ProductionRepository interface {
	GetDefProductions() ([]model.Production, error)
	CreateDefProduction(production model.Production) error
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
		return nil, fmt.Errorf("could not query def producion rows: %w", err)
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
	return nil
}
