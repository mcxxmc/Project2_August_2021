package dbInterface

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"serverGo/src/common"
	"strconv"
)

// For a template "fake" record; it does not have a real image of that name; used only for test.
var templatePNG = "template.PNG"
var templatePath = "D:/Project2_August_2021/s3/temp/template.PNG"

// SQL commands.
var sqlDropTable = "DROP TABLE IF EXISTS picture;"
var sqlCreateTable = "CREATE TABLE picture" + "(" +
	"id INTEGER NOT NULL AUTO_INCREMENT," +
	"name VARCHAR(100) UNIQUE NOT NULL," +
	"path VARCHAR(200) NOT NULL," +
	"prediction BOOLEAN," +
	"label BOOLEAN," +
	"PRIMARY KEY(id)," +
	"UNIQUE (name)" + ");"
var sqlInsert = "INSERT INTO picture(name, path, prediction, label) values(?,?,?,?);"
var sqlInsertWithPrediction = "INSERT INTO picture(name, path, prediction) values(?,?,?);"
var sqlInsertWithLabel = "INSERT INTO picture(name, path, label) values(?,?,?);"
var sqlInsertBared = "INSERT INTO picture(name, path) values(?,?);"
var sqlUpdatePrediction = "UPDATE picture SET prediction=? WHERE name=?;"
var sqlUpdatePathAndPrediction = "UPDATE picture SET path=?, prediction=? WHERE name=?"
var sqlUpdateLabel = "UPDATE picture SET label=? WHERE name=?;"
var sqlUpdatePathAndLabel = "UPDATE picture SET path=?, label=? WHERE name=?;"
var sqlQueryTry = "SELECT name FROM picture WHERE id=1;"
var sqlQueryName = "SELECT path, prediction, label FROM picture WHERE name=?;"
var sqlFetchAll = "SELECT id, name, path, prediction, label FROM picture;"
var sqlLimit = " LIMIT "
var comma = ","
var sqlFetchUnlabeled = "SELECT name, path FROM picture WHERE label IS NULL"
var sqlFetchUnpredictedUnlabeled = "SELECT name, path FROM picture WHERE prediction IS NULL AND label IS NULL;"

// OpenDb opens a connection to the database.
func OpenDb() *sql.DB {
	// connect to a MySQL server as a user named "localuser" to a database named august2021
	db, err := sql.Open("mysql", "localuser:localuserpassword@tcp(localhost:3306)/august2021")
	common.CheckErr(err)
	fmt.Println("Db connection opens.")
	return db
}

// CloseDb closes the connection.
func CloseDb(db *sql.DB) {
	err := db.Close()
	common.CheckErr(err)
	fmt.Println("Db connection closes.")
}

// Tests if the database functions properly.
func testDb(db *sql.DB) error {
	_, err := db.Query(sqlQueryTry)
	common.CheckErr(err)
	return err
}

// RecreateTable deletes the old table and rebuild a new one with default values.
func RecreateTable(db *sql.DB) {
	_, err := db.Exec(sqlDropTable)
	common.CheckErr(err)

	_, err = db.Exec(sqlCreateTable)
	common.CheckErr(err)

	_, err = db.Exec(sqlInsert, templatePNG, templatePath, nil, true)
	common.CheckErr(err)
}

// TryConnection tries the connection to the database and checks if it functions properly.
func TryConnection() {
	// TODO: change the query to insert then delete template.PNG; there should not be any fake data
	db := OpenDb()
	defer CloseDb(db)
	// Open doesn't open a connection. Check if the connection is valid and stop the program if it is not.
	err := db.Ping()
	common.PanicErr(err)
	// Check if the table exist:
	err = testDb(db)
	if err != nil {
		RecreateTable(db)
		fmt.Println("Table rebuilt.")
	}else{
		fmt.Println("The database is ready.")
	}
}

// Insert inserts a new record into the database.
func Insert(name string, path string, prediction bool, label bool) {
	db := OpenDb()
	defer CloseDb(db)
	_, err := db.Exec(sqlInsert, name, path, prediction, label)
	common.CheckErr(err)

	// id, err := res.LastInsertId()
	// common.CheckErr(err)
	// fmt.Println("The last inserted Id is: ", id)
}

// InsertWithPrediction inserts a new record into the database.
func InsertWithPrediction(name string, path string, prediction bool) {
	db := OpenDb()
	defer CloseDb(db)
	_, err := db.Exec(sqlInsertWithPrediction, name, path, prediction)
	common.CheckErr(err)
}

// InsertWithLabel inserts a new record into the database.
func InsertWithLabel(name string, path string, label bool) {
	db := OpenDb()
	defer CloseDb(db)
	_, err := db.Exec(sqlInsertWithLabel, name, path, label)
	common.CheckErr(err)
}

// InsertBared inserts a new record with name and path into the database.
func InsertBared(name string, path string) {
	db := OpenDb()
	defer CloseDb(db)
	_, err := db.Exec(sqlInsertBared, name, path)
	common.CheckErr(err)
}

// UpdatePrediction updates the prediction attribute.
func UpdatePrediction(name string, prediction bool) {
	db := OpenDb()
	defer CloseDb(db)
	_, err := db.Exec(sqlUpdatePrediction, prediction, name)
	common.CheckErr(err)

	// n, err := res.RowsAffected()
	// common.CheckErr(err)
	// fmt.Printf("Updates Succeeds. Affected rows: %d\n", n)
}

// UpdatePathAndPrediction updates the path and the prediction.
func UpdatePathAndPrediction(name string, path string, prediction bool) {
	db := OpenDb()
	defer CloseDb(db)
	_, err := db.Exec(sqlUpdatePathAndPrediction, path, prediction, name)
	common.CheckErr(err)
}

// UpdateLabel updates the prediction attribute.
func UpdateLabel(name string, label bool) {
	db := OpenDb()
	defer CloseDb(db)
	_, err := db.Exec(sqlUpdateLabel, label, name)
	common.CheckErr(err)
}

// UpdatePathAndLabel updates the path and the label.
func UpdatePathAndLabel(name string, path string, label bool) {
	db := OpenDb()
	defer CloseDb(db)
	_, err := db.Exec(sqlUpdatePathAndLabel, path, label, name)
	common.CheckErr(err)
}

// QueryName checks if the name is in the database.
// Returns a tuple, of which the first bool means whether the name is in the database.
// The second bool is "prediction" and the third bool is "label".
// If the record with the name does not exist, the 2 booleans will always be nil.
// The last element in the tuple is a string, which is  the path attribute of the record,
// and it will be nil if the record does not exist.
func QueryName(name string) (bool, *bool, *bool, *string) {
	db := OpenDb()
	defer CloseDb(db)
	res, err := db.Query(sqlQueryName, name)

	defer func(res *sql.Rows) {
		err := res.Close()
		common.CheckErr(err)
	}(res)

	common.CheckErr(err)

	var path string
	var prediction *bool  // should be declared as pointers
	var label *bool
	if res.Next(){
		err := res.Scan(&path, &prediction, &label)
		common.CheckErr(err)
		return true, prediction, label, &path
	}else {
		return false, nil, nil, nil
	}
}

// FetchAll fetches all the records and returns id, name, path, prediction and label.
func FetchAll() Records {
	db := OpenDb()
	defer CloseDb(db)

	res, err := db.Query(sqlFetchAll)

	defer func(res *sql.Rows) {
		err := res.Close()
		common.CheckErr(err)
	}(res)

	common.CheckErr(err)

	var records Records  // initialize a slice first
	var id int
	var name string
	var path string
	var prediction *bool
	var label *bool

	for {
		if res.Next(){
			err := res.Scan(&id, &name, &path, &prediction, &label)
			common.CheckErr(err)
			records.Recs = append(records.Recs,
				Record{Id: id, Name: name, Path: path, Prediction: prediction, Label: label})
		} else{
			break
		}
	}

	return records
}

// FetchN fetches the first n records starting from the offset.
func FetchN(offset int, n int) []PathAndDesc {
	var r []PathAndDesc
	var id int
	var name string
	var path string
	var prediction *bool
	var label *bool
	command := sqlFetchAll + sqlLimit + strconv.Itoa(offset) + comma + strconv.Itoa(n)
	db := OpenDb()
	defer CloseDb(db)
	res, err := db.Query(command)
	defer func(res *sql.Rows) {
		err := res.Close()
		common.CheckErr(err)
	}(res)
	common.CheckErr(err)
	for {
		if res.Next(){
			err := res.Scan(&id, &name, &path, &prediction, &label)
			common.CheckErr(err)
			r = append(r, PathAndDesc{Path: path, Text: generateText(id, name, prediction, label)})
		}else {
			break
		}
	}
	return r
}

// FetchUnlabeled fetches all the unlabeled records and return the names and the paths.
func FetchUnlabeled() []UnlabeledRecord {
	var r []UnlabeledRecord
	var name string
	var path string
	db := OpenDb()
	defer CloseDb(db)
	res, err := db.Query(sqlFetchUnlabeled)
	defer func(res *sql.Rows) {
		err := res.Close()
		common.CheckErr(err)
	}(res)
	common.CheckErr(err)
	for {
		if res.Next(){
			err := res.Scan(&name, &path)
			common.CheckErr(err)
			r = append(r, UnlabeledRecord{Name: name, Path: path})
		}else {
			break
		}
	}
	return r
}

// FetchUnpredictedUnlabeled fetches all the records that are neither predicted nor labeled.
func FetchUnpredictedUnlabeled() []UnpredictedUnlabeledRecord {
	var r []UnpredictedUnlabeledRecord
	var name string
	var path string
	db := OpenDb()
	defer CloseDb(db)
	res, err := db.Query(sqlFetchUnpredictedUnlabeled)
	defer func(res *sql.Rows) {
		err := res.Close()
		common.CheckErr(err)
	}(res)
	common.CheckErr(err)
	for {
		if res.Next() {
			err := res.Scan(&name, &path)
			common.CheckErr(err)
			r = append(r, UnpredictedUnlabeledRecord{Name: name, Path: path})
		}else {
			break
		}
	}
	return r
}
