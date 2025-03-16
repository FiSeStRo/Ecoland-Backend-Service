package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
)

type ProductRepository interface {
	GetDefProducts() ([]model.Product, error)
	CreateProduct(product model.Product) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetDefProducts() ([]model.Product, error) {

	defProductQuery := `SELECT * FROM def_product`

	rows, err := r.db.Query(defProductQuery)
	if err != nil {
		return nil, fmt.Errorf("could not select def products: %w", err)
	}
	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Value); err != nil {
			return nil, fmt.Errorf("can not read from product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) CreateProduct(product model.Product) error {
	defProductQuery := `INSERT INTO def_product(token_name, base_value) VALUES(?,?)`

	result, err := r.db.Exec(defProductQuery, product.Name, product.Value)

	if err != nil {
		return fmt.Errorf("failed to insert new product: %w", err)
	}

	if id, err := result.LastInsertId(); id <= 0 || err != nil {
		return fmt.Errorf("could not insert new product: %w", err)
	}
	return nil
}
