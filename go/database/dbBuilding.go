package database

func GetPossibleProductionsOfBuilding(buildingId int) ([]DefProduction, error) {

	var listOfProductions []DefProduction
	rows, err := db.Query(`SELECT rel.production_id, p.token_name, p.cost, p.base_duration FROM def_rel_building_production rel JOIN def_production p ON p.id = rel.production_id WHERE rel.building_id=? `, buildingId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var production DefProduction
		err = rows.Scan(&production.Id, &production.TokenName, &production.Cost, &production.BaseDuration)
		if err != nil {
			return nil, err
		}
		listOfProductions = append(listOfProductions, production)
	}

	return listOfProductions, nil

}

func FindBuilding(id int) (Building, error) {

	var building Building
	row := db.QueryRow(`SELECT user_id, def_id, name, time_build FROM buildings WHERE id=? `, id)

	err := row.Scan(&building.UserId, &building.DefId, &building.Name, &building.TimeBuild)
	if err != nil {
		return Building{}, err
	}
	building.Id = id
	return building, nil
}

func GetBuildingStorage(id int) ([]BuildingStorage, error) {
	var storage []BuildingStorage
	rows, err := db.Query(`SELECT building_id, product_id, amount, capacity FROM rel_building_product WHERE building_id=? `, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := BuildingStorage{}
		err = rows.Scan(&item.BuildingId, &item.ProductId, &item.Capacity, &item.Amount)
		if err != nil {
			return nil, err
		}
		storage = append(storage, item)
	}

	return storage, nil
}
