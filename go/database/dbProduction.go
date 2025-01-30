package database

func FindProductionById(id int) (Production, error) {

	var production Production
	row := db.QueryRow(`SELECT * FROM ? WHERE id=?`, RelBuildingDefProductionTable, id)
	err := row.Scan(&production.Id, &production.BuildingId, &production.ProductionId, &production.TimeStart, &production.TimeEnd, &production.IsCompleted)
	if err != nil {
		return Production{}, nil
	}
	return production, nil
}
