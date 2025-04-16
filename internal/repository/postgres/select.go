package postgres

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNoRecordFound = errors.New("No record found")

func SelectAllRecords(db *gorm.DB, orderBy string, models interface{}, receiver interface{}) error {

	if orderBy == "" {
		orderBy = "asc"
	}
	res := db.Model(models).Order(orderBy).Find(receiver)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return ErrNoRecordFound
	}
	return nil
}

func SelectById(db *gorm.DB, id string, models interface{}, receiver interface{}, ids ...string) error {
	res := db.Model(models).Where(`id = ?`, id).Find(receiver)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return ErrNoRecordFound
	}
	return nil
}

// TODO: handle select with pagination and groupby

func CheckExistsAndReturnModelInstance(db *gorm.DB, query string, models interface{}, args ...interface{}) (interface{}, error) {
	err := db.Where(query, args...).First(models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func CheckExists(db *gorm.DB, query string, models interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := db.Model(models).Where(query, args...).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func Preload(db *gorm.DB, query, preload string, model interface{}, args ...interface{}) error {
	err := db.Model(model).Preload(preload, query, args).First(model).Error
	if err != nil {
		return err
	}
	return nil
}
