package db

import (
	"golang.org/x/oauth2"

	_ "github.com/mattn/go-sqlite3"
)

type AuthProfile struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type User struct {
	Profile   AuthProfile
	AuthToken oauth2.Token
}

type DB struct {
	Users []User
}

var database = DB{}

func Init() {
	database = DB{}
}

func AddUser(user User) {
	database.Users = append(database.Users, user)
}

func GetUsers() []User {
	return database.Users
}
