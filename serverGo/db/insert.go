package db

import (
	"gorm.io/gorm"
	"serverGo/common"
)

// Insert inserts a new record into the database.
func Insert(db *gorm.DB, name string, path string, prediction bool, label bool) {
	res := db.Create(&RecordNoId{
		Name: name,
		Path: path,
		Prediction: &prediction,
		Label: &label,
	})
	common.CheckErr(res.Error)
}

// InsertWithPrediction inserts a new record into the database.
func InsertWithPrediction(db *gorm.DB, name string, path string, prediction bool) {
	res := db.Create(&RecordNoId{
		Name: name,
		Path: path,
		Prediction: &prediction,
	})
	common.CheckErr(res.Error)
}

// InsertWithLabel inserts a new record into the database.
func InsertWithLabel(db *gorm.DB, name string, path string, label bool) {
	res := db.Create(&RecordNoId{
		Name: name,
		Path: path,
		Label: &label,
	})
	common.CheckErr(res.Error)
}

// InsertBared inserts a new record with name and path into the database.
func InsertBared(db *gorm.DB, name string, path string) {
	res := db.Create(&RecordNoId{Name: name, Path: path})
	common.CheckErr(res.Error)
}
