package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(dbURL string) *pgxpool.Pool {
	dbUrlWithDisableSSL := fmt.Sprintf("%s?sslmode=disable", dbURL)
	pool, err := pgxpool.New(context.Background(), dbUrlWithDisableSSL)
	if err != nil {
		log.Fatal("Cannot connect to the Database: ", err)
	}
	return pool
}
