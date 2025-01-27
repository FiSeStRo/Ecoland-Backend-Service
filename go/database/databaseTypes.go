package database

type DefBuildingTable struct {
	Id                   int     `json:"id"`
	TokenName            string  `json:"token_name"`
	BaseConstructionCost float64 `json:"base_construction_cost"`
	BaseConstructionTime int     `json:"base_construction_time"`
}

type DefProductTable struct {
	Id        int    `json:"id"`
	BaseValue int    `json:"base_value"`
	TokenName string `json:"token_name"`
}

type DefProductionTable struct {
	Id           int     `json:"id"`
	TokenName    string  `json:"token_name"`
	Cost         float64 `json:"cost"`
	BaseDuration int     `json:"base_duration"`
}
