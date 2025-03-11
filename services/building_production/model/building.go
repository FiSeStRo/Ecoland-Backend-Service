package model

// Building represents a structure that can be built
type Building struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	ResourceCost float64 `json:"resourceCost"`
	BuildTime    int     `json:"buildTime"` // in minutes
	Production   []int   `json:"production"`
}

// GetAllBuildings returns all available building types
func GetAllBuildings() []Building {

	return []Building{
		{ID: 1, Name: "Farm", Description: "Produces food", ResourceCost: 1000, BuildTime: 30, Production: []int{1, 2, 3}},
		{ID: 2, Name: "Mine", Description: "Produces ore", ResourceCost: 2000, BuildTime: 60, Production: []int{1, 2}},
		{ID: 3, Name: "Factory", Description: "Processes raw materials", ResourceCost: 5000, BuildTime: 120, Production: []int{3}},
	}
}
