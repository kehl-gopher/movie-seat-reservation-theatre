package postgres

import (
	"gorm.io/gorm"
)

func DeleteAllRecords(db *gorm.DB, model interface{}) error {
	err := db.Delete(model).Error
	return err
}

func DeleteWithPrimaryKey(db *gorm.DB, model interface{}, ids ...string) error {
	err := db.Delete(model, ids).Error
	return err
}

func DeleteSingleRecord(db *gorm.DB, query string, model interface{}, args ...interface{}) error {
	err := db.Delete(model, query, args).Error
	return err
}

func ForceDelete(db *gorm.DB, query string, model interface{}, args ...interface{}) error {
	err := db.Unscoped().Where(query, args).Delete(model).Error
	return err
}
