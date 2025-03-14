package database

type DefBuilding struct {
	ID          int    `json:"id"`
	Name        string `json:"token_name"`
	Cost        int    `json:"base_construction_cost"`
	BuildTime   int    `json:"base_construction_time"`
	Productions []int  `json:"productions"`
}

type DefProduction struct {
	ID       int     `json:"id"`
	Name     string  `json:"token_name"`
	Cost     float64 `json:"cost"`
	Duration int     `json:"base_duration"`
}

type DefRelProductionProduct struct {
	ProductionID int  `json:"production_id"`
	ProductID    int  `json:"product_id"`
	IsInput      bool `json:"is_input"`
	Amount       int  `json:"amount"`
}
