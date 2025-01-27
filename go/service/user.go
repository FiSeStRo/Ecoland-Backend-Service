package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/FiSeStRo/Ecoland-Backend-Service/authentication"
	"github.com/FiSeStRo/Ecoland-Backend-Service/database"
	"github.com/FiSeStRo/Ecoland-Backend-Service/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserJson struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	Id       int
	Username string
	Password []byte
	Email    string
}

type SignInResponse struct {
	Id           int    `json:"id"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// SignIn hadles the user SignIn
func SignIn(w http.ResponseWriter, req *http.Request) {

	if !utils.IsMethodPOST(w, req) {
		return
	}

	var signInUser UserJson
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "unable to ready body", 400)
	}
	err = json.Unmarshal(body, &signInUser)
	if err != nil {
		http.Error(w, "Unable to read body", 400)
	}

	ok, err := UserExists(signInUser.Username)
	if err != nil || !ok {
		http.Error(w, "Wrong username or password", 400)
		return
	}

	user, err := FindUser(signInUser.Username)
	if err != nil {
		http.Error(w, "Wrong username or password", 400)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(signInUser.Password))
	if err != nil {
		http.Error(w, "Wrong username or password", 400)
	}

	accessToken, err := authentication.CreateNewJWT(user.Id, false)
	if err != nil {
		http.Error(w, "ups something went wrong", 500)
		return
	}
	refreshToken, err := authentication.CreateNewJWT(user.Id, true)
	if err != nil {
		http.Error(w, "ups something went wrong", 500)
		return
	}
	resBody := SignInResponse{
		Id:           user.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resBody)
}

// SignUp allows for users to sign up for the service
func SignUp(w http.ResponseWriter, req *http.Request) {
	var signUpUser UserJson

	if !utils.IsMethodPOST(w, req) {
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "unalbe to read the body", 400)
		return
	}

	err = json.Unmarshal(body, &signUpUser)
	if err != nil {
		http.Error(w, "could not decode JSON", 400)
		return
	}

	ue, err := UserExists(signUpUser.Username)
	if err != nil {
		http.Error(w, "couldn not search for Username", 400)
		return
	}

	if ue {
		http.Error(w, "Username alreday exists", 400)
		return
	}

	if !utils.IsEmail(signUpUser.Email) {
		http.Error(w, "In valid email", http.StatusBadRequest)
	}

	bp, err := bcrypt.GenerateFromPassword([]byte(signUpUser.Password), bcrypt.MinCost)

	if err != nil {
		http.Error(w, "upps something went wrong please try again", 400)
		return
	}

	newUser := User{
		Username: signUpUser.Username,
		Password: bp,
		Email:    signUpUser.Email,
	}

	err = CreateUser(newUser)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not add new user", 400)
		return
	}

	w.WriteHeader(201)
}

type UserRequestBody struct {
	Id int `json:"userId"`
}

type UserResourcesResBody struct {
	Money    float64 `json:"money"`
	Prestige int     `json:"prestige"`
}

func GetUserResources(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "WrongMethod", http.StatusMethodNotAllowed)
		return
	}

	claims, err := authentication.ValidateAuthentication(req)

	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
	}

	db := database.GetDB()

	stmt, err := db.Prepare(`SELECT money, prestige FROM user_resources WHERE user_id = ?`)
	if err != nil {
		http.Error(w, "Could not prepare db query", 500)
		return
	}

	var resBody UserResourcesResBody
	err = stmt.QueryRow(claims.UserId).Scan(&resBody.Money, &resBody.Prestige)
	if err != nil {
		http.Error(w, "Could not find User resources", 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resBody)
}
func AddUser(user User) (int, error) {
	db := database.GetDB()
	stmt, err := db.Prepare(`INSERT into users(username, password) VALUES (?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Password)
	if err != nil {
		return 0, err
	}
	userId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(userId), err
}

func UserExists(username string) (bool, error) {
	db := database.GetDB()
	stmt, err := db.Prepare(`SELECT COUNT(*) FROM users WHERE username = ?`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func FindUser(username string) (User, error) {
	db := database.GetDB()
	stmt, err := db.Prepare(`SELECT id, username, password FROM users WHERE username = ?`)

	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	var user User
	err = stmt.QueryRow(username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil

}

// CreateUser creates a user with user details and inits user resources
func CreateUser(user User) error {

	userId, err := AddUser(user)
	if err != nil {
		return err
	}
	err = CreateNewUserResources(userId)
	if err != nil {
		return err
	}

	return nil
}

func CreateNewUserResources(userId int) error {

	db := database.GetDB()
	initMoney := 100000
	initPrestige := 10
	stmt, err := db.Prepare(`INSERT into user_resources(user_id , money, prestige) VALUES(?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, initMoney, initPrestige)
	if err != nil {
		return err
	}
	return nil
}
