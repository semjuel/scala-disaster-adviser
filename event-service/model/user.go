package model

import (
	"scala-disaster-adviser/event-service/database"
)

type User struct {
	Uuid  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func UserSave(user User) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	exists := UserFindOneByEmail(user.Email)
	if exists == nil {
		return nil
	}

	insertQuery := "INSERT INTO users (uuid, name, email, token) VALUES ($1, $2, $3, $4)"
	err = db.QueryRow(insertQuery, user.Uuid, user.Name, user.Email, user.Token).Scan()

	return err
}

func UserFindOneByEmail(email string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	var id int = 0
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if id != 0 {
		return nil
	}

	return err
}
