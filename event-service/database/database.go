package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"sync"
)

type Database struct {
	Instance *sql.DB
	err      error
}

var once sync.Once

func (db *Database) Connect() error {
	once.Do(func() {
		dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"),
		)

		connect, err := sql.Open("postgres", dbInfo)
		if err != nil {
			db.err = err
		} else {
			db.err = nil
			db.Instance = connect
		}
	})

	return db.err
}
