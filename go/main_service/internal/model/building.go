package model

import (
	"time"
)

// BuildingDefinition represents a type of building that can be constructed
type BuildingDefinition struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	BaseCost float64 `json:"baseCost"`
	BaseTime int     `json:"baseTime"` // in minutes
	CanBuild bool    `json:"canBuild"` // calculated field, not stored in DB
}

// Building represents a player-owned building instance
type Building struct {
	ID                    int       `json:"id"`
	UserID                int       `json:"userId"`
	BuildingDefID         int       `json:"buildingDefId"`
	Status                string    `json:"status"` // "under_construction" or "operational"
	ConstructionStartTime time.Time `json:"constructionStartTime,omitempty"`
	ConstructionEndTime   time.Time `json:"constructionEndTime,omitempty"`
}

// Storage represents products stored in a building
type Storage struct {
	ID           int    `json:"id"`
	BuildingID   int    `json:"buildingId"`
	ProductDefID int    `json:"productDefId"`
	Quantity     int    `json:"quantity"`
	ProductName  string `json:"productName"` // from joined def_products table
}

// Production represents a production configuration in a building
type Production struct {
	ID           int     `json:"id"`
	BuildingID   int     `json:"buildingId"`
	ProductDefID int     `json:"productDefId"`
	Status       string  `json:"status"`
	Rate         float64 `json:"rate"`
	ProductName  string  `json:"productName"` // from joined def_products table
}

// ConstructRequest represents a request to construct a new building
type ConstructRequest struct {
	BuildingDefID int `json:"buildingDefId"`
}

// ConstructResult represents the result of a building construction
type ConstructResult struct {
	BuildingID     int       `json:"buildingId"`
	CompletionTime time.Time `json:"completionTime"`
	RemainingFunds float64   `json:"remainingFunds"`
}

// BuildingDetails combines building data with its related information
type BuildingDetails struct {
	Building   Building           `json:"building"`
	Definition BuildingDefinition `json:"definition"`
	Storage    []Storage          `json:"storage,omitempty"`
	Production []Production       `json:"production,omitempty"`
}
