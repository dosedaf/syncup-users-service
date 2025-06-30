package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() (*pgx.Conn, error) {
	fmt.Println("Starting connection with Postgres Db")
	connStr := "postgres://postgres:password@localhost:5432/db?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = conn.Ping(context.Background()); err != nil {
		log.Println("DB Ping Failed")
		log.Fatal(err)
	}

	log.Println("DB Connection started successfully")

	return conn, nil
}
