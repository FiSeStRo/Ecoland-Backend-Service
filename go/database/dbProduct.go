package database

type ProductionProduct struct {
	DefProduct
	IsInput bool `json:"is_input"`
	Amount  int  `json:"amount"`
}

func FindProductsByProductionId(id int) ([]ProductionProduct, error) {
	var products []ProductionProduct
	rows, err := db.Query(`SELECT p.id, p.token_name, p.base_value, rel.is_input, rel.amount FROM def_rel_production_product rel JOIN def_product p ON p.id = rel.product_id WHERE rel.production_id=?`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var product ProductionProduct
		err = rows.Scan(&product.Id, &product.TokenName, &product.BaseValue, &product.IsInput, &product.Amount)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
