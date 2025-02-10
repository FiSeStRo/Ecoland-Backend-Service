package database

func FindUserById(userId int) (User, error) {

	stmt, err := db.Prepare(`SELECT id, username, password FROM users WHERE id = ?`)

	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	var user User
	err = stmt.QueryRow(userId).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil

}

func GetUserResources(userId int) (UserResource, error) {
	var userResource UserResource
	row := db.QueryRow(`SELECT user_id, money, prestige FROM user_resources WHERE user_id = ?`, userId)
	err := row.Scan(&userResource.UserId, &userResource.Money, &userResource.Prestige)
	if err != nil {
		return UserResource{}, err
	}
	return userResource, nil
}

type PUserResource struct {
	UserId   int
	Money    *float64
	Prestige *int
}

func UpdateUserResources(res PUserResource) error {
	if res.Money != nil {
		_, err := db.Exec(`INSERT INTO ?(money) VALUES(?) WHERE id=?`, UserResourceTable, *res.Money, res.UserId)
		if err != nil {
			return err
		}
	} else if res.Prestige != nil {
		_, err := db.Exec(`INSERT INTO ?(prestige) VALUES(?) WHERE id=?`, UserResourceTable, *res.Prestige, res.UserId)
		if err != nil {
			return err
		}
	}
	return nil
}

func EmailExists(email string) (bool, error) {
	rslt := db.QueryRow(`SELECT COUNT(*) FROM users WHERE email=?`, email)
	var count int
	err := rslt.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateUserInfo(user User) error {
	_, err := db.Exec(`INSERT INTO ?(username, password, email) WHERE id=?`, user.Username, user.Password, user.Email, user.Id)
	if err != nil {
		return err
	}
	return nil
}
