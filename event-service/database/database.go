package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

// Connect to the database.
func Connect() (*sql.DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	return sql.Open("postgres", dbInfo)
}
