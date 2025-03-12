package model

// Building represents a structure that can be built
type Building struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	ResourceCost float64 `json:"resourceCost"`
	BuildTime    int     `json:"buildTime"` // in minutes
	Productions  []int   `json:"productions"`
}
