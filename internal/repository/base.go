package repository

import "gorm.io/gorm"

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(db *gorm.DB) *Postgres {
	return &Postgres{
		DB: db,
	}
}

type Database struct {
	Pdb *Postgres
}

var DB = &Database{}

func ConnectDB() *Database {
	return DB
}
