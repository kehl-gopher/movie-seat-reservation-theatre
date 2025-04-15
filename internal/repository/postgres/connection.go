package postgres

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(config *env.Config) (*gorm.DB, error) {
	// verify connection port

	port, err := verifyPort(config.DBPORT)

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

func verifyPort(port string) (int, error) {
	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, errors.New("invalid port expected integer")
	}

	if p < 0 || p > 65535 {
		fmt.Println("Port provided is an invalid postgres port.... falling back to default postgres port")
		p = 5432
	}

	return p, nil
}
