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
	Page      uint   `json:"page"`
	Limit     uint   `json:"limit"`
	TotalPage uint   `json:"total_page"`
	PrevPage  string `json:"prev_page,omitempty"`
	NextPage  string `json:"next_page,omitempty"`
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

func SelectAllRecords(db *gorm.DB, orderBy, order string, models interface{}, receiver interface{}) error {

	if orderBy == "" && order == "" {
		orderBy = fmt.Sprintf("%s %s", "id", "asc")
	} else {
		orderBy = fmt.Sprintf("%s %s", order, orderBy)
	}
	res := db.Model(models).Order(orderBy).Find(receiver)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return ErrNoRecordFound
	}
	return nil
}

func buildPrevNextPage(pagination Pagination, totalPage uint) (string, string) {
	if pagination.Offset > totalPage {
		return "", ""
	}
	if pagination.Offset == 1 {
		if totalPage == 1 {
			return "", ""
		}
		nextPage := fmt.Sprintf("movies?page=%d&limit=%d", pagination.Offset+1, pagination.Limit)
		prevPage := ""
		return nextPage, prevPage
	} else if pagination.Offset == totalPage {
		if totalPage == 1 {
			return "", ""
		}
		nextPage := ""
		prevPage := fmt.Sprintf("movies?page=%d&limit=%d", pagination.Offset-1, pagination.Limit)
		return nextPage, prevPage
	} else if pagination.Offset < totalPage {
		nextPage := fmt.Sprintf("movies?page=%d&limit=%d", pagination.Offset+1, pagination.Limit)
		prevPage := fmt.Sprintf("movies?page=%d&limit=%d", pagination.Offset-1, pagination.Limit)
		return nextPage, prevPage
	}
	nextPage := fmt.Sprintf("movies?page=%d&limit=%d", pagination.Offset+1, pagination.Limit)
	prevPage := fmt.Sprintf("movies?page=%d&limit=%d", pagination.Offset-1, pagination.Limit)
	return nextPage, prevPage
}
func SelectAllRecordWithPagination(db *gorm.DB, query string, orderBy, order string, model interface{}, receiver interface{}, limit, offset uint, args ...interface{}) (PaginationResponse, error) {
	var count int64

	if orderBy == "" {
		orderBy = "asc"
	}
	if order == "" {
		order = "id"
	}

	pagination := GetPagination(offset, limit)

	err := db.Model(model).Count(&count).Error
	if err != nil {
		return PaginationResponse{
			Page:  pagination.Offset,
			Limit: pagination.Limit,
		}, err
	}
	off := int(pagination.Limit) * (int(pagination.Offset) - 1)
	tx := db.Model(model).Limit(int(pagination.Limit)).Offset((off)).Order(fmt.Sprintf("%s %s", order, orderBy)).Find(receiver)
	if tx.Error != nil {
		return PaginationResponse{
			Page:  pagination.Offset,
			Limit: pagination.Limit,
		}, tx.Error
	}
	nextPage, prevPage := buildPrevNextPage(pagination, uint(count)/pagination.Limit)
	return PaginationResponse{
		Page:      pagination.Offset,
		Limit:     pagination.Limit,
		TotalPage: uint(count) / pagination.Limit,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}, nil
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
		if errors.Is(gorm.ErrRecordNotFound, res.Error) {
			return ErrNoRecordFound
		}
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
