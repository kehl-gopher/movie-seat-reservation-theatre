package main

import (
	"fmt"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models/migration"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models/seeding"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/minio"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/redis"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/router"
)

func main() {
	config := env.SetEnv()
	//connect to database
	_, err := postgres.ConnectPostgres(config)
	if err != nil {
		fmt.Println("failed to connect to postgres database", err)
		panic(err)
	}
	_, err = redis.ConnectRedis(config)
	if err != nil {
		fmt.Println("failed to connect to redis database", err)
	}
	_, err = minio.ConnectMinio(config)
	if err != nil {
		fmt.Println("failed to connect to minio database", err)
		panic(err)
	}
	db := repository.ConnectDB()

	// run migration for models
	if err := migration.RunMigrations(db); err == nil {
		seeding.StartSeeding(db)
	} else {
		panic(err)
	}
	router.Router(config, db)
}
