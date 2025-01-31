package database

func FindProductionById(id int) (Production, error) {

	var production Production
	row := db.QueryRow(`SELECT * FROM ? WHERE id=?`, BuildingProductionTable, id)
	err := row.Scan(&production.Id, &production.BuildingId, &production.ProductionId, &production.TimeStart, &production.TimeEnd, &production.IsCompleted)
	if err != nil {
		return Production{}, nil
	}
	return production, nil
}

func FindProductionByBuilding(buildingId int) ([]Production, error) {
	productions := make([]Production, 0)
	rows, err := db.Query(`SELECT * FROM ? WHERE building_id=?`, BuildingProductionTable, buildingId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var production Production
		err = rows.Scan(&production)
		if err != nil {
			return nil, err
		}
		productions = append(productions, production)
	}

	return productions, nil
}
