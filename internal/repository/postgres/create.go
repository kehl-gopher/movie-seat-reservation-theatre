package postgres

import (
	"errors"

	"gorm.io/gorm"
)

func Create(db *gorm.DB, model interface{}) error {
	result := db.Create(model)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("No data insert occured")
	}

	return nil
}

func CreateMany[T any](db *gorm.DB, model []T) error {
	result := db.Create(model)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != int64(len(model)) {
		return errors.New("incomplete data insert occured")
	}

	return nil
}
