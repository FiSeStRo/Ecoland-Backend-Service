package database

import "fmt"

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

func GetBStorage(id int) ([]BStorage, error) {
	var storageList []BStorage
	rows, err := db.Query(`SELECT * FROM ? WHERE building_id=?`, StorageBuildingTable, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var storage BStorage
		err := rows.Scan(&storage.BuildingId, &storage.ProductId, &storage.Amount, &storage.Capacity)
		if err != nil {
			return nil, err
		}
		storageList = append(storageList, storage)
	}

	return storageList, nil
}

func GetBStorageOfP(buildingId int, productId int) (BStorage, error) {
	var storage BStorage
	row := db.QueryRow(`SELECT * FROM ? WHERE building_id=? AND product_id=?`, StorageBuildingTable, buildingId, productId)
	err := row.Scan(&storage.BuildingId, &storage.ProductId, &storage.Amount, &storage.Capacity)
	if err != nil {
		return BStorage{}, err
	}
	return storage, nil
}

func UpdateBStorageOfP(storage BStorage) error {
	rslt, err := db.Exec(`UPDATE ? SET amount=?, capacity=? WHERE building_id =? AND product_id=?`, StorageBuildingTable, storage.Amount, storage.Capacity, storage.BuildingId, storage.ProductId)
	if err != nil {
		return err
	}
	rows, err := rslt.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("failed to update product %v in building %v", storage.ProductId, storage.BuildingId)
	}
	return nil
}

func BatchUpdateBStorageOfP(storages []BStorage) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, storage := range storages {
		rslt, err := tx.Exec(`UPDATE ? SET amount=?, capacity=? WHERE building_id =? AND product_id=?`, StorageBuildingTable, storage.Amount, storage.Capacity, storage.BuildingId, storage.ProductId)
		if err != nil {
			return err
		}
		rows, err := rslt.RowsAffected()
		if err != nil || rows == 0 {
			return fmt.Errorf("failed to update product %v in building %v", storage.ProductId, storage.BuildingId)
		}
		return nil
	}
	return tx.Commit()
}

func IsUserBuildingOwner(userId int, buildingId int) bool {
	row := db.QueryRow(`SELECT user_id FROM ? WHERE building_id=?`, BuildingsTable, buildingId)
	var id int
	if err := row.Scan(&id); err != nil || id != userId {
		return false
	}
	return true
}
