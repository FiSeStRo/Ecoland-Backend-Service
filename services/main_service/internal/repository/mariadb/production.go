package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/repository"
)

type productionRepository struct {
	db *sql.DB
}

func NewProductionRepository(db *sql.DB) repository.ProductionRepository {
	return &productionRepository{db: db}
}

func (r *productionRepository) GetProductionDefenitions() ([]model.ProductionDefenition, error) {
	query := `SELECT id, token_name, cost, base_duration FROM def_production`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var definitions []model.ProductionDefenition
	for rows.Next() {
		var def model.ProductionDefenition
		err := rows.Scan(&def.ID, &def.TokenName, &def.Cost, &def.BaseDuration)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		definitions = append(definitions, def)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return definitions, nil
}

func (r *productionRepository) GetProductsOfProduction(productionID int) ([]model.ProductOfProduction, error) {
	stmt, err := r.db.Prepare(`SELECT p.product_id, p.is_input, p.amount, def.token_name 
			FROM def_rel_production_product p 
			JOIN def_product def ON p.prodcut_id = def.id 
			WHERE p.production_id=?`)

	if err != nil {
		return nil, fmt.Errorf("could not prepare db: %w", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(productionID)

	var baseProducts []model.ProductOfProduction

	for rows.Next() {
		var bP model.ProductOfProduction

		err := rows.Scan(&bP.Id, &bP.IsInput, &bP.Amount, &bP.TokenName)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		baseProducts = append(baseProducts, bP)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}

	return baseProducts, nil
}

func (r *productionRepository) GetBuildingIdOfProduction(productionID int) (int, error) {
	query := `SELECT building_id FROM rel_building_def_production WHERE id=?`
	row := r.db.QueryRow(query, productionID)
	var buildingID int
	if err := row.Scan(&buildingID); err != nil {
		return 0, fmt.Errorf("error finding buildingID: %w", err)
	}
	return buildingID, nil
}

func (r *productionRepository) DeleteProductionByID(productionID int) error {
	query := `DELETE FROM rel_buildng_def_production WHERE id=?`
	row, err := r.db.Exec(query, productionID)
	if err != nil {
		return fmt.Errorf("error deleting from db: %w", err)
	}
	if affected, err := row.RowsAffected(); err != nil || affected == 0 {
		return fmt.Errorf("no rows where effected: %w", err)
	}
	return nil
}

func (r *productionRepository) GetProductionsByBuildingType(buildingDefID int) ([]model.ProductionDefenition, error) {
	query := `SELECT rel.production_id, p.token_name, p.cost, p.base_duration FROM def_rel_building_production rel JOIN def_production p ON p.id = rel.production_id WHERE rel.building_id=? `
	rows, err := r.db.Query(query, buildingDefID)
	if err != nil {
		return nil, fmt.Errorf("error finding productions: %w", err)
	}
	var productions []model.ProductionDefenition
	for rows.Next() {
		var production model.ProductionDefenition
		if err := rows.Scan(&production.ID, &production.TokenName, &production.Cost, &production.BaseDuration); err != nil {
			return nil, fmt.Errorf("error scanning producitons: %w", err)
		}
		productions = append(productions, production)
	}
	return productions, nil
}

func (r *productionRepository) AddProduction(newProduction model.BaseProduction) error {
	exec := `INSERT INTO rel_building_def_production(building_id, production_id, time_start, time_end, cycles, is_completed) VALUES(?,?,?,?,?,?)`
	stmt, err := r.db.Prepare(exec)
	if err != nil {
		return fmt.Errorf("error preparing db: %w", err)
	}
	defer stmt.Close()

	rslt, err := stmt.Exec(exec, newProduction.BuildingID, newProduction.ProductionID, newProduction.TimeStart, newProduction.TimeEnd, newProduction.Cycles, false)
	if err != nil {
		return fmt.Errorf("error adding new production: %w", err)
	}

	if id, err := rslt.LastInsertId(); id <= 0 || err != nil {
		return fmt.Errorf("no production added: %w", err)
	}

	return nil
}
