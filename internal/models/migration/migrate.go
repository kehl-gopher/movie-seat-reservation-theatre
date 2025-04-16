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
		&models.UserRoleID{},
		&models.AccessToken{},
	}
}

func RunMigrations(db *repository.Database) error {
	return AutoMigrate(db.Pdb.DB, MigrateModel()...)

}

func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Model migration successful")
	return nil
}
