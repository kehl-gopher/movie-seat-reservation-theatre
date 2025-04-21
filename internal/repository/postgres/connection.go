package postgres

import (
	"fmt"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(config *env.Config) (*gorm.DB, error) {
	port, err := utility.VerifyPort(config.DBPORT)

	if err != nil {
		return nil, err
	}

	db, err := connectionPostgres(config.DBNAME, config.DBHOST, config.DBPASSWORD, config.DBUSER, config.DB_SSLMODE, config.DB_TIMEZONE, port)

	if err != nil {
		return nil, err
	}

	repository.DB.Pdb = repository.NewPostgres(db)

	fmt.Println("Database connection successful")
	return db, nil
}

func connectionPostgres(dbName, host, password, user, sslmode, timezone string, port int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, user, password, dbName, port, sslmode, timezone)
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("failed to connect to postgres database")
		return nil, err
	}

	return db, nil
}
