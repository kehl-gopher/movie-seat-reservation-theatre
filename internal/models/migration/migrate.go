package migration

import (
	"fmt"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"gorm.io/gorm"
)

func MigrateModel() []interface{} {
	return []interface{}{
		&models.Permission{},
		&models.Role{},
		&models.Users{},
	}
}

func RunMigrations(db *repository.Database) {
	AutoMigrate(db.Pdb.DB, MigrateModel()...)
}

func AutoMigrate(db *gorm.DB, models ...interface{}) {
	err := db.AutoMigrate(models...)
	if err != nil {
		panic(err)
	}

	fmt.Println("Model migration successful")
}
