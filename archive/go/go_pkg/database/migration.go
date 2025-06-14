package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/FiSeStRo/Ecoland-Backend-Service/go_pkg/config"
)

// Migration represents a database schema change
type Migration struct {
	Name string
	SQL  string
}

// BuildingTableMigrations returns migrations for building-related tables
func BuildingTableMigrations() []Migration {
	const (
		defBuildingsTableSQL = `CREATE TABLE IF NOT EXISTS def_buildings(
            id INT PRIMARY KEY AUTO_INCREMENT,
            token_name VARCHAR(255),
            base_construction_cost DECIMAL(10,2),
            base_construction_time INT
        )`

		buildingsTableSQL = `CREATE TABLE IF NOT EXISTS buildings(
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            def_id INT NOT NULL,
            name VARCHAR(255),
            lan INT NOT NULL,
            lat INT NOT NULL,
            time_build TIMESTAMP NOT NULL
        )`
	)

	var DefTable = Migration{
		Name: "Create def_buildings table",
		SQL:  defBuildingsTableSQL,
	}

	var BuildingsTable = Migration{
		Name: "Create buildings table",
		SQL:  buildingsTableSQL,
	}

	return []Migration{
		DefTable,
		BuildingsTable,
	}
}

// UserTableMigrations returns migrations for user-related tables
func UserTableMigrations() []Migration {
	const (
		userTableSQL = `CREATE TABLE IF NOT EXISTS users(
                id INT AUTO_INCREMENT PRIMARY KEY,
                username VARCHAR(255) UNIQUE NOT NULL,
                password VARCHAR(255) NOT NULL,
                email VARCHAR(255) NOT NULL,
                role INT,
                time_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                time_last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )`
		UserResourceTableSQL = `CREATE TABLE IF NOT EXISTS user_resources(
                user_id int PRIMARY KEY,
                money decimal(10,2) NOT NULL DEFAULT 100000,
                prestige int
            )`
	)

	var UserTable = Migration{Name: "Create users table",
		SQL: userTableSQL}

	var UserResourceTable = Migration{
		Name: "Create user_resource table",
		SQL:  UserResourceTableSQL,
	}
	return []Migration{
		UserTable,
		UserResourceTable,
	}
}

// ProductTableMigrations returns migrations for product-related tables
func ProductTableMigrations() []Migration {
	const defProductTableSQL = `CREATE TABLE IF NOT EXISTS def_product(
        id INT PRIMARY KEY AUTO_INCREMENT,
        base_value INT NOT NULL,
        token_name VARCHAR(255)
    )`

	var DefProductTable = Migration{
		Name: "Create def_product table",
		SQL:  defProductTableSQL,
	}

	return []Migration{
		DefProductTable,
	}
}

// ProductionTableMigrations returns migrations for production-related tables
func ProductionTableMigrations() []Migration {
	const defProductionTableSQL = `CREATE TABLE IF NOT EXISTS def_production(
        id INT PRIMARY KEY AUTO_INCREMENT,
        token_name VARCHAR(255),
        cost DECIMAL(10,2),
        base_duration INT
    )`

	var DefProductionTable = Migration{
		Name: "Create def_production table",
		SQL:  defProductionTableSQL,
	}

	const defRelProductionProductSQL = `CREATE TABLE IF NOT EXISTS def_rel_production_product(
        production_id INT NOT NULL,
        product_id INT NOT NULL,
        is_input BOOLEAN NOT NULL COMMENT 'True if consumed, False if produced',
        amount INT NOT NULL,
        PRIMARY KEY (production_id, product_id, is_input),
        FOREIGN KEY (production_id) REFERENCES def_production(id),
        FOREIGN KEY (product_id) REFERENCES def_product(id)
    )`

	var DefRelProductionProductTable = Migration{
		Name: "Create def_rel_production_product table",
		SQL:  defRelProductionProductSQL,
	}
	return []Migration{
		DefProductionTable,
		DefRelProductionProductTable,
	}
}

// RelationTableMigrations returns migrations for relation tables
func RelationTableMigrations() []Migration {
	const (
		relBuildingProductTableSQL = `CREATE TABLE IF NOT EXISTS rel_building_product(
            building_id INT NOT NULL,
            product_id INT NOT NULL,
            amount INT,
            capacity INT,
            PRIMARY KEY (building_id, product_id)
        )`

		defRelBuildingProductionTableSQL = `CREATE TABLE IF NOT EXISTS def_rel_building_production(
            building_id INT NOT NULL,
            production_id INT NOT NULL,
            PRIMARY KEY (building_id, production_id)
        )`

		relBuildingDefProductionTableSQL = `CREATE TABLE IF NOT EXISTS rel_building_def_production(
            id INT AUTO_INCREMENT PRIMARY KEY,
            building_id INT NOT NULL,
            production_id INT NOT NULL,
            time_start TIMESTAMP NOT NULL,
            time_end TIMESTAMP NOT NULL,
            cycles INT NOT NULL,
            is_completed BOOL DEFAULT FALSE
        )`
	)

	var RelBuildingProductTable = Migration{
		Name: "Create rel_building_product table",
		SQL:  relBuildingProductTableSQL,
	}

	var DefRelBuildingProductionTable = Migration{
		Name: "Create def_rel_building_production table",
		SQL:  defRelBuildingProductionTableSQL,
	}

	var RelBuildingDefProductionTable = Migration{
		Name: "Create rel_building_def_production table",
		SQL:  relBuildingDefProductionTableSQL,
	}

	return []Migration{
		RelBuildingProductTable,
		DefRelBuildingProductionTable,
		RelBuildingDefProductionTable,
	}
}

// MapTableMigrations returns migrations for map-related tables
func MapTableMigrations() []Migration {
	const (
		cityTableSQL = `CREATE TABLE IF NOT EXISTS city(
            lon INT,
            lat INT,
            name VARCHAR(255),
            pop INT NOT NULL,
            PRIMARY KEY (lon, lat)
        )`

		productMapTableSQL = `CREATE TABLE IF NOT EXISTS product_map(
            lon INT,
            lat INT,
            product_id INT,
            PRIMARY KEY (lon, lat, product_id)
        )`
	)

	var CityTable = Migration{
		Name: "Create city table",
		SQL:  cityTableSQL,
	}

	var ProductMapTable = Migration{
		Name: "Create product_map table",
		SQL:  productMapTableSQL,
	}

	return []Migration{
		CityTable,
		ProductMapTable,
	}
}

// StorageTableMigrations returns migrations for storage-related tables
func StorageTableMigrations() []Migration {
	const storageTableSQL = `CREATE TABLE IF NOT EXISTS storage(
        id INT AUTO_INCREMENT PRIMARY KEY,
        building_id INT NOT NULL,
        product_def_id INT NOT NULL,
        quantity DECIMAL(10,2) NOT NULL DEFAULT 0,
        UNIQUE KEY (building_id, product_def_id)
    )`

	var StorageTable = Migration{
		Name: "Create storage table",
		SQL:  storageTableSQL,
	}

	return []Migration{
		StorageTable,
	}
}

// ProductionEntryTableMigrations returns migrations for production entries
func ProductionEntryTableMigrations() []Migration {
	const productionTableSQL = `CREATE TABLE IF NOT EXISTS production(
        id INT AUTO_INCREMENT PRIMARY KEY,
        building_id INT NOT NULL,
        product_def_id INT NOT NULL,
        status VARCHAR(50) NOT NULL DEFAULT 'inactive',
        production_rate DECIMAL(10,2) NOT NULL DEFAULT 0,
        UNIQUE KEY (building_id, product_def_id)
    )`

	var ProductionTable = Migration{
		Name: "Create production table",
		SQL:  productionTableSQL,
	}

	return []Migration{
		ProductionTable,
	}
}

// RunMigrations executes all database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []Migration{}
	migrations = append(migrations, BuildingTableMigrations()...)
	migrations = append(migrations, UserTableMigrations()...)
	migrations = append(migrations, ProductTableMigrations()...)
	migrations = append(migrations, ProductionTableMigrations()...)
	migrations = append(migrations, RelationTableMigrations()...)
	migrations = append(migrations, MapTableMigrations()...)
	migrations = append(migrations, StorageTableMigrations()...)
	migrations = append(migrations, ProductionEntryTableMigrations()...)

	for _, migration := range migrations {
		if _, err := db.Exec(migration.SQL); err != nil {
			return fmt.Errorf("migration '%s' failed: %w", migration.Name, err)
		}
		fmt.Printf("Migration completed: %s\n", migration.Name)
	}

	return nil
}

// MigrateDefinitonData migrates and sets up the definition tables in the database
func MigrateDefinitonData(db *sql.DB) error {
	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return fmt.Errorf("could not disable foreign key checks: %w", err)
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS def_buildings"); err != nil {
		return fmt.Errorf("could not drop def_buildings table: %w", err)
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS def_production"); err != nil {
		return fmt.Errorf("could not drop def_production table: %w", err)
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS def_product"); err != nil {
		return fmt.Errorf("could not drop def_product table: %w", err)
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS def_rel_building_production"); err != nil {
		return fmt.Errorf("could not drop def_rel_building_production: %w", err)
	}
	if _, err := db.Exec("DROP TABLE IF EXISTS def_rel_production_product"); err != nil {
		return fmt.Errorf("could not drop def_rel_production_product: %w", err)
	}
	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=1"); err != nil {
		return fmt.Errorf("could not enable foreign key checks: %w", err)
	}

	if err := RunMigrations(db); err != nil {
		return fmt.Errorf("Migration failed: %w", err)
	}

	log.Println("Start DefData Migration ...")
	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return fmt.Errorf("could not disable foreign key checks: %w", err)
	}
	defBuildings, err := config.LoadJsonDataFromFileStorage[[]DefBuilding]("def_buildings.json")
	if err != nil {
		return fmt.Errorf("couldn't load json from file: %w", err)
	}

	for _, building := range defBuildings {
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
		buildingQuery := `INSERT INTO def_buildings(token_name, base_construction_cost, base_construction_time) VALUES ( ?, ?, ?)`
		buildingResult, err := tx.Exec(buildingQuery, building.Name, building.Cost, building.BuildTime)
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
	}

	type DefProductionConfig struct {
		DefProduction
		Products []struct {
			ProductID int  `json:"product_id"`
			IsInput   bool `json:"is_input"`
			Amount    int  `json:"amount"`
		}
	}

	defProductions, err := config.LoadJsonDataFromFileStorage[[]DefProductionConfig]("def_production.json")

	if err != nil {
		return fmt.Errorf("could not load json from file %w", err)
	}

	for _, production := range defProductions {
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("faield to begin transaction: %w", err)
		}

		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		productionQuery := `INSERT INTO def_production( token_name, cost, base_duration) VALUES(?, ?, ?)`
		result, err := tx.Exec(productionQuery, production.Name, production.Cost, production.Duration)
		if err != nil {
			return fmt.Errorf("failed to insert production: %w", err)
		}

		productionID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last inserted ID: %w", err)
		}
		if len(production.Products) > 0 {
			productQuery := `INSERT INTO def_rel_production_product(production_id, product_id, is_input, amount) VALUE(?, ?, ?, ?)`
			stmt, err := tx.Prepare(productQuery)
			if err != nil {
				return fmt.Errorf("failed to prepare relation statement: %w", err)
			}
			defer stmt.Close()
			for _, product := range production.Products {
				if _, err := stmt.Exec(productionID, product.ProductID, product.IsInput, product.Amount); err != nil {
					return fmt.Errorf("failed to insert products production relation: %w", err)
				}
			}
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	defProducts, err := config.LoadJsonDataFromFileStorage[[]DefProduct]("def_products.json")
	if err != nil {
		return fmt.Errorf("Could not load def_product.json : %w", err)
	}
	productQuery := `INSERT INTO def_product(token_name, base_value) VALUE(?, ?)`
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("faield to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, product := range defProducts {
		_, err := tx.Exec(productQuery, product.Name, product.Value)
		if err != nil {
			return fmt.Errorf("failed to insert products %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=1"); err != nil {
		return fmt.Errorf("could not enable foreign key checks: %w", err)
	}
	log.Println("Migrate Definiton Data Successful")
	return nil
}
