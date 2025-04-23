package postgres

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrNoRecordFound = errors.New("No record found")

var (
	defaultPageNumaber uint = 1
	defaultPageSize    uint = 10
)

type Pagination struct {
	Offset uint
	Limit  uint
}

type PaginationResponse struct {
	Page      uint `json:"page"`
	Limit     uint `json:"limit"`
	TotalPage uint `json:"total_page"`
}

func GetPagination(offset, limit uint) Pagination {
	if offset == 0 {
		offset = defaultPageNumaber
	}
	if limit == 0 {
		limit = defaultPageSize
	}
	return Pagination{
		Offset: offset,
		Limit:  limit,
	}
}

func SelectAllRecords(db *gorm.DB, orderBy, value string, models interface{}, receiver interface{}) error {

	if orderBy == "" && value == "" {
		orderBy = fmt.Sprintf("%s %s", "id", "asc")
	} else {
		orderBy = fmt.Sprintf("%s %s", value, orderBy)
	}
	res := db.Model(models).Order(orderBy).Find(receiver)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return ErrNoRecordFound
	}
	return nil
}

func SelectAllRecordWithPagination(db *gorm.DB, query string, models interface{}, receiver interface{}, limit, offset uint, args ...interface{}) {
	
}
func SelectMultipleRecord(db *gorm.DB, query string, model interface{}, receiver interface{}, args ...interface{}) error {
	res := db.Model(model).Where(query, args...).Scan(receiver)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
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
		return false, err
	}
	return count > 0, nil
}

func Preload(db *gorm.DB, model interface{}, preloads ...string) *gorm.DB {
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}
