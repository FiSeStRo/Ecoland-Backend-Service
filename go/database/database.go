package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// SetDB sets the used Database
func SetDB(database *sql.DB) {
	db = database
}

// GetDB gets the used Database
func GetDB() *sql.DB {
	return db
}

// InitDatabaseTables initalises necessary DatabaseTables
func InitDatabaseTables() {
	initaliseBuildingTables()
	initialiseUserTable()
	initaliseProductTables()
	initaliseProductionTables()
	initaliseRelationTables()
}
func initaliseBuildingTables() {

	createDefBuildingsTable := `CREATE TABLE IF NOT EXISTS def_buildings(
	id INT PRIMARY KEY,
	token_name VARCHAR(255),
	base_construction_cost DECIMAL,
	base_construction_time INT
	)`

	_, err := db.Exec(createDefBuildingsTable)
	if err != nil {
		log.Panic("error creating def_building_table", err)
	} else {
		log.Println("def_buildings table created or already in place")
	}

	createBuildingsTable := `CREATE TABLE IF NOT EXISTS buildings(
	id INT AUTO_INCREMENT PRIMARY KEY,
	user_id INT NOT NULL,
	def_id INT NOT NULL,
	name VARCHAR(255),
	lan INT NOT NULL,
	lat INT NOT NULL,
	time_build TIMESTAMP NOT NULL
	)`

	_, err = db.Exec(createBuildingsTable)
	if err != nil {
		log.Panic("error creating buildings_table", err)
	} else {
		log.Println("buildings table created or already in place")
	}
}

func initialiseUserTable() {
	createUserTable := `CREATE TABLE IF NOT EXISTS users(
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		role INT,
		time_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		time_last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	createUserResourcesTable := `CREATE TABLE IF NOT EXISTS user_resources(
  		user_id int PRIMARY KEY,
  		money decimal NOT NULL DEFAULT 100000,
  		prestige int
		);`

	_, err := db.Exec(createUserTable)
	if err != nil {
		log.Panic("error creating user Table", err)
	} else {
		log.Println("userTable created or already in place")
	}
	_, err = db.Exec(createUserResourcesTable)
	if err != nil {
		log.Panic("error creating user_resources Table", err)
	} else {
		log.Println("user_resources Table created or already in place")
	}
}

func initaliseProductTables() {
	createDefProducts := `CREATE TABLE IF NOT EXISTS def_product(
	id INT PRIMARY KEY,
	base_value INT NOT NULL,
	token_name VARCHAR(255))`

	_, err := db.Exec(createDefProducts)
	if err != nil {
		log.Panic("error creating def_product table", err)
	} else {
		log.Println("def_product Table created or already in place")
	}

}

func initaliseProductionTables() {

	createDefProduction := `CREATE TABLE IF NOT EXISTS def_production(
	id INT PRIMARY KEY,
	token_name VARCHAR(255),
	COST DECIMAL,
	BASE_DURATION INT
	)`

	_, err := db.Exec(createDefProduction)
	if err != nil {
		log.Panic("error creating def_production table", err)
	} else {
		log.Println("def_production Table created or already in place")
	}

}

func initaliseRelationTables() {
	createRelBuildingProduct := `CREATE TABLE IF NOT EXISTS rel_building_product(
	building_id INT NOT NULL,
	product_id INT NOT NULL,
	amount INT,
	capacity INT
	)`

	_, err := db.Exec(createRelBuildingProduct)
	if err != nil {
		log.Panic("error creating def_product table", err)
	} else {
		log.Println("def_product Table created or already in place")
	}

	createDefRelBuildingProduction := `CREATE TABLE IF NOT EXISTS def_rel_building_production(
	building_id int NOT NULL,
	production_id int NOT NULL
	)`

	_, err = db.Exec(createDefRelBuildingProduction)
	if err != nil {
		log.Panic("error creating def_rel_building_production table", err)
	} else {
		log.Println("def_rel_building_production Table created or already in place")
	}

	createRelBuildingDefProduction := `CREATE TABLE IF NOT EXISTS rel_buildng_def_production(
	id INT  AUTO_INCREMENT PRIMARY KEY,
	building_id INT NOT NULL,
	production_id INT NOT NULL,
	time_start TIMESTAMP NOT NULL,
	time_end TIMESTAMP NOT NULL,
	cycles INT NOT NULL,
	is_completed BOOL DEFAULT FALSE
	)`

	_, err = db.Exec(createRelBuildingDefProduction)
	if err != nil {
		log.Panic("error creating rel_buildng_def_production table", err)
	} else {
		log.Println("rel_buildng_def_production Table created or already in place")
	}
}
