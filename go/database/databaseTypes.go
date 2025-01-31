package database

import "time"

type DefBuilding struct {
	Id                   int     `json:"id"`
	TokenName            string  `json:"token_name"`
	BaseConstructionCost float64 `json:"base_construction_cost"`
	BaseConstructionTime int     `json:"base_construction_time"`
}

type DefProduct struct {
	Id        int    `json:"id"`
	BaseValue int    `json:"base_value"`
	TokenName string `json:"token_name"`
}

type DefProduction struct {
	Id           int     `json:"id"`
	TokenName    string  `json:"token_name"`
	Cost         float64 `json:"cost"`
	BaseDuration int     `json:"base_duration"`
}

type Building struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	DefId     int       `json:"def_id"`
	Name      string    `json:"name"`
	TimeBuild time.Time `json:"time_build"`
}

type BuildingStorage struct {
	BuildingId int `json:"building_id"`
	ProductId  int `json:"product_id"`
	Amount     int `json:"amount"`
	Capacity   int `json:"capacity"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password []byte `json:"password"`
	Email    string `json:"email"`
}

type UserResource struct {
	UserId   int     `json:"user_id"`
	Money    float64 `json:"money"`
	Prestige int     `json:"prestige"`
}

type Production struct {
	Id           int
	BuildingId   int
	ProductionId int
	TimeStart    int
	TimeEnd      int
	Cycles       int
	IsCompleted  bool
}

type BStorage struct {
	BuildingId int `json:"building_id"`
	ProductId  int `json:"product_id"`
	Amount     int `json:"amount"`
	Capacity   int `json:"capacity"`
}

const RelBuildingDefProductionTable = "rel_building_def_production"
const UserResourceTable = "user_resources"
const StorageBuildingTable = "rel_building_product"
const BuildingsTable = "buildings"
