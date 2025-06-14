package model

import "time"

type BaseProduct struct {
	Id        int    `json:"id"`
	TokenName string `json:"token_name"`
	Amount    int    `json:"amount"`
}
type ProductionModel struct {
	Id           int           `json:"id"`
	TokenName    string        `json:"token_name"`
	BaseDuration time.Time     `json:"base_duration"`
	Cost         float64       `json:"cost"`
	ProductsIn   []BaseProduct `json:"products_in"`
	ProductsOut  []BaseProduct `json:"products_out"`
}

type ProductionDefenition struct {
	ID           int     `json:"id"`
	TokenName    string  `json:"token_name"`
	Cost         float64 `json:"cost"`
	BaseDuration int64   `json:"base_duration"`
}

type ProductOfProduction struct {
	BaseProduct
	IsInput bool `json:"is_input"`
}

type NewProductionRequest struct {
	ProductionID int `json:"id"`
	BuildingID   int `json:"building_id"`
	Cycles       int `json:"cycles"`
}

type BaseProduction struct {
	ProductionID int
	BuildingID   int
	TimeStart    int64
	TimeEnd      int64
	Cycles       int
}

type ProductionTable struct {
	BaseProduction
	ID          int
	IsCompleted bool
}
