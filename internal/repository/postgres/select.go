package postgres

import (
	"errors"
	"fmt"

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

func SelectSingleRecord(db *gorm.DB, query string, model interface{}, receiver interface{}, args ...interface{}) error {
	err := db.Model(model).Where(query, args).First(receiver).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return ErrNoRecordFound
		}
		return err
	}
	return nil
}

func SelectById(db *gorm.DB, id interface{}, models interface{}, receiver interface{}) error {
	res := db.Model(models).Where(`id = ?`, id).First(receiver)
	if res.Error != nil {
		return res.Error
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
		fmt.Println("------------>")
		return false, err
	}
	return count > 0, nil
}

func Preload(db *gorm.DB, preload string, model interface{}, args ...interface{}) error {
	res := db.Model(model).Preload(preload, args).Find(model)
	if res.Error != nil {
		fmt.Println(res.Error, "---------->")
		return res.Error
	} else if res.RowsAffected == 0 {
		return errors.New("No record found")
	}
	return nil
}
