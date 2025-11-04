package db

import (
	"context"
	"fmt"
	"log"

	"github.com/Stenoliv/didlydoodash_api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToDatabase() (*pgxpool.Pool, error) {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	conn, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		fmt.Println("Cannot connect to database! ")
		log.Fatal("Connection error: ", err)
	}
	fmt.Println("Connected to the database!")

	return conn, nil
}

func Load() (*pgxpool.Pool, error) {
	db, err := ConnectToDatabase()
	if err != nil {
		return nil, err
	}
	return db, nil
}
