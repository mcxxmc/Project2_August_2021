package db

import (
	"gorm.io/gorm"
	"webserver/common"
)

// UpdatePrediction updates the prediction attribute.
func UpdatePrediction(db *gorm.DB, name string, prediction bool) {
	res := db.Model(&Record{}).Where(conditionNameEqual, name).Update(attrPred, prediction)
	common.CheckErr(res.Error)
}

// UpdatePathAndPrediction updates the path and the prediction.
func UpdatePathAndPrediction(db *gorm.DB, name string, path string, prediction bool) {
	res := db.Model(&Record{}).Where(conditionNameEqual, name).Updates(Record{Path: path, Prediction: &prediction})
	common.CheckErr(res.Error)
}

// UpdateLabel updates the prediction attribute.
func UpdateLabel(db *gorm.DB, name string, label bool) {
	res := db.Model(&Record{}).Where(conditionNameEqual, name).Update(attrLabel, label)
	common.CheckErr(res.Error)
}

// UpdatePathAndLabel updates the path and the label.
func UpdatePathAndLabel(db *gorm.DB, name string, path string, label bool) {
	res := db.Model(&Record{}).Where(conditionNameEqual, name).Updates(Record{Path: path, Label: &label})
	common.CheckErr(res.Error)
}
