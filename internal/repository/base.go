package repository

import (
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(db *gorm.DB) *Postgres {
	return &Postgres{
		DB: db,
	}
}

type Rediss struct {
	Rdb *redis.Client
}

func NewRedis(db *redis.Client) *Rediss {
	return &Rediss{
		Rdb: db,
	}
}

type Database struct {
	Pdb *Postgres
	Red *Rediss
	Min *minio.Client
}

var DB = &Database{}

func ConnectDB() *Database {
	return DB
}
