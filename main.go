package main

import (
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models/migration"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models/seeding"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository/postgres"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/router"
)

func main() {

	config := env.SetEnv()

	//connect to database
	_, err := postgres.ConnectPostgres(config)
	if err != nil {
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
