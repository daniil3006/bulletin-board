package main

import (
	"bulletin-board/internal/ad/repository/pgstore"
	"bulletin-board/internal/transport/api"
	"bulletin-board/pkg/postgresql"
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	pc := postgresql.PostgresConfig{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
	}

	pool, err := postgresql.NewClient(ctx, pc)
	if err != nil {
		log.Fatal("error to connect to PostgreSQL")
	}
	defer pool.Close()

	log.Println("Success connect to PostgreSQL!")

	store := pgstore.NewRepository(pool)
	r := api.NewRouter(store)
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
