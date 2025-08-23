package main

import (
	"bulletin-board/internal/ad/repository/pgstore"
	"bulletin-board/internal/ad/service"
	"bulletin-board/internal/ad/transport/api"
	userPgstore "bulletin-board/internal/user/pgstore"
	userServ "bulletin-board/internal/user/service"
	userApi "bulletin-board/internal/user/transport/api"
	"bulletin-board/pkg/postgresql"
	"context"
	"github.com/gorilla/mux"
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

	adRepo := pgstore.NewRepository(pool)
	adService := service.NewService(adRepo)
	adHandler := api.NewHandler(*adService)

	userRepo := userPgstore.NewRepository(pool)
	userService := userServ.NewService(userRepo)
	userHandler := userApi.NewHandler(*userService)

	r := mux.NewRouter()
	adHandler.NewRouter(r)
	userHandler.NewRouter(r)

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Route: %s Methods: %v", path, methods)
		return nil
	})

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
