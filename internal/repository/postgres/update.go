package postgres

import "gorm.io/gorm"

func UpdateSingleRecord(db *gorm.DB, query string, model interface{}, args ...interface{}) error {
	tx := db.Model(model).Where(query, args).Updates(model)
	return tx.Error
}

func UpdateRelationShipRecord(db *gorm.DB, model interface{}, newAssociation map[string]interface{}) error {
	var err error
	for key, value := range newAssociation {
		err = db.Model(model).Association(key).Replace(value)
		if err != nil {
			return err
		}
	}

	return nil
}
