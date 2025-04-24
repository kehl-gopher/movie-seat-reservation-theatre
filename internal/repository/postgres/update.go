package postgres

import "gorm.io/gorm"




func UpdateSingleRecord(db *gorm.DB, query string, model interface{}, args ...interface{}) error {
	err := db.Model(model).Where(query, args).Updates(model).Error
	return err
}
