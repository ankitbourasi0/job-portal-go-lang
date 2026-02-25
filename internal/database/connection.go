package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Postgres driver
)

func InitDB(dbURL string) *sql.DB {
	dbUrlWithDisableSSL := fmt.Sprintf("%s?sslmode=disable", dbURL)
	db, err := sql.Open("postgres", dbUrlWithDisableSSL)
	if err != nil {
		log.Fatal("Cannot connect to the Database: ", err)
	}

	//check connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Database Ping failed: ", err)
	}

	return db
}
