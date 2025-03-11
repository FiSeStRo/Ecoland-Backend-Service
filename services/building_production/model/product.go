package model

// Product represents items that can be produced
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Value       float64 `json:"value"`
	Weight      float64 `json:"weight"`
	Category    string  `json:"category"`
}

// GetAllProducts returns all available products
func GetAllProducts() []Product {
	// In a real application, this would fetch from database
	return []Product{
		{ID: 1, Name: "Wheat", Description: "Basic crop", Value: 5, Weight: 0.5, Category: "Food"},
		{ID: 2, Name: "Iron Ore", Description: "Raw mineral", Value: 10, Weight: 2.0, Category: "Material"},
		{ID: 3, Name: "Tool", Description: "Crafted implement", Value: 20, Weight: 1.0, Category: "Equipment"},
	}
}
