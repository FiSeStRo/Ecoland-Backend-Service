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
