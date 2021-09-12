package db

var dataSource = "localuser:localuserpassword@tcp(localhost:3306)/august2021"
var tableName = "picture"

var conditionNameEqual = "name = ?"
var conditionIdBiggerThan = "id > ?"
var conditionLabelIsNull = "label is null"
var conditionPredAndLabelAreNull = "prediction is null and label is null"
// var attrName = "name"
// var attrPath = "path"
var attrPred = "prediction"
var attrLabel = "label"

// RecordNoId used to insert new record into the database
type RecordNoId struct {
	Name string  `gorm:"column:name"`
	Path string  `gorm:"column:path"`
	Prediction *bool  `gorm:"column:prediction"`
	Label *bool  `gorm:"column:label"`
}

// TableName manually set the table name for this struct
func (RecordNoId) TableName() string {
	return tableName
}

// Record A complete record (including id, name, path, prediction, label) in the database.
type Record struct {
	Id int `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
	Path string `json:"path" gorm:"column:path"`
	Prediction *bool `json:"prediction" gorm:"column:prediction"`
	Label *bool `json:"label" gorm:"column:label"`
}

func (Record) TableName() string {
	return tableName
}

// Records A collection of Record s.
type Records struct {
	Recs []Record `json:"records"`
}

func (Records) TableName() string {
	return tableName
}

/*
// PathAndDesc The JSON structure containing the path of an image and a customized description of that image.
type PathAndDesc struct {
	Path string
	Text string
}

func (PathAndDesc) TableName() string {
	return tableName
}

// UnlabeledRecord The JSON structure containing the information (name and path) of an unlabeled image.
type UnlabeledRecord struct {
	Name string
	Path string
}

type UnpredictedUnlabeledRecord struct {
	Name string
	Path string
}

// Close the sql rows.
func closeSqlRows(res *sql.Rows) {
	err := res.Close()
	common.CheckErr(err)
}

 */
