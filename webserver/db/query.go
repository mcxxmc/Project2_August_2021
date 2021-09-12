package db

import (
	"gorm.io/gorm"
	"log"
	"webserver/common"
)

// QueryName checks if the name is in the database. Returns a tuple with 4 values.
// first: bool. Whether the name is in the database.
// second: *bool. "Prediction". nil if the record does not exist.
// third: *bool. "Label". nil if the record does not exist.
// forth: *string. "Path". nil if the record does not exist.
func QueryName(db *gorm.DB, name string) (bool, *bool, *bool, *string) {
	var user Record
	res := db.Where(conditionNameEqual, name).Take(&user)
	if res.Error == nil {
		return true, user.Prediction, user.Label, &user.Path
	} else {
		// log.Println(res.Error)
		return false, nil, nil, nil
	}
}

// FetchAll fetches all the records and returns id, name, path, prediction and label.
func FetchAll(db *gorm.DB) Records {
	var records Records
	res := db.Find(&records.Recs)
	common.CheckErr(res.Error)
	log.Println("db/query.go FetchAll: number of records fetched: ", res.RowsAffected)
	return records
}

// FetchN fetches the first n records starting from the offset.
func FetchN(db *gorm.DB, offset int, n int) Records {
	var records Records
	res := db.Where(conditionIdBiggerThan, offset).Limit(n).Find(&records.Recs)
	common.CheckErr(res.Error)
	log.Println("db/query.go FetchN: number of records fetched: ", res.RowsAffected)
	return records
}

// FetchUnlabeled fetches all the unlabeled records and return the names and the paths.
func FetchUnlabeled(db *gorm.DB) Records {
	var records Records
	res := db.Where(conditionLabelIsNull).Find(&records.Recs)
	common.CheckErr(res.Error)
	log.Println("db/query.go FetchUnlabeled: number of records fetched: ", res.RowsAffected)
	return records
}

// FetchUnpredictedUnlabeled fetches all the records that are neither predicted nor labeled.
func FetchUnpredictedUnlabeled(db *gorm.DB) Records {
	var records Records
	res := db.Where(conditionPredAndLabelAreNull).Find(&records.Recs)
	common.CheckErr(res.Error)
	log.Println("db/query.go FetchUnpredictedUnlabeled: number of records fetched: ", res.RowsAffected)
	return records
}
