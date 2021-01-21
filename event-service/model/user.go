package model

import (
	"log"
)

type User struct {
	Id    int64  `json:"id"`
	Uuid  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func UserSave(user User) error {
	err := DB.Connect()
	if err != nil {
		return err
	}

	exists := UserFindOneByEmail(user.Email)
	if exists == nil {
		return nil
	}

	insertQuery := "INSERT INTO users (uuid, name, email, token) VALUES ($1, $2, $3, $4)"
	err = DB.Instance.QueryRow(insertQuery, user.Uuid, user.Name, user.Email, user.Token).Scan()

	return err
}

func UserFindOneByEmail(email string) error {
	err := DB.Connect()
	if err != nil {
		return err
	}

	var id int = 0
	err = DB.Instance.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if id != 0 {
		return nil
	}

	return err
}

func UserAll() ([]User, error) {
	var users []User

	err := DB.Connect()
	if err != nil {
		return users, err
	}

	rows, err := DB.Instance.Query("SELECT u.id, u.uuid, u.name, u.email, u.token FROM users u")
	if err != nil {
		log.Printf("Wrong request to the db, %v", err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Token)
		if err == nil {
			users = append(users, user)
		} else {
			log.Printf("Error %s", err)
		}
	}

	return users, nil
}
