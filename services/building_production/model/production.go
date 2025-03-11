package model

// Production represents a production process
type Production struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	InputType   string  `json:"inputType"`
	OutputType  string  `json:"outputType"`
	Duration    int     `json:"duration"` // in minutes
	Cost        float64 `json:"cost"`
}

// GetAllProductions returns all available production processes
func GetAllProductions() []Production {
	// In a real application, this would fetch from database
	return []Production{
		{ID: 1, Name: "Basic Farming", Description: "Simple crop cultivation", InputType: "Seeds", OutputType: "Crops", Duration: 60, Cost: 100},
		{ID: 2, Name: "Mining Operation", Description: "Extract minerals from mine", InputType: "Tools", OutputType: "Raw Ore", Duration: 120, Cost: 200},
		{ID: 3, Name: "Manufacturing", Description: "Process raw materials into products", InputType: "Raw Materials", OutputType: "Finished Goods", Duration: 180, Cost: 500},
	}
}
