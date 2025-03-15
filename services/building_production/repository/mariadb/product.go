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
	return nil
}
