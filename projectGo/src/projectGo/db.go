package projectGo

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)


var templatePNG = "template.PNG"
var templatePath = "D:/Project2_August_2021/s3/temp/template.PNG"

var sqlDropTable = "DROP TABLE IF EXISTS picture;"

var sqlCreateTable = "CREATE TABLE picture" + "(" +
	"id INTEGER NOT NULL AUTO_INCREMENT," +
	"name VARCHAR(100) UNIQUE NOT NULL," +
	"path VARCHAR(200) NOT NULL," +
	"prediction BOOLEAN," +
	"label BOOLEAN," +
	"PRIMARY KEY(id)," +
	"UNIQUE (name)" + ");"

var sqlInsert = "INSERT INTO picture(name, path, prediction, label) values(?,?,?,?)"
var sqlInsertWithPrediction = "INSERT INTO picture(name, path, prediction) values(?,?,?,)"
var sqlInsertWithLabel = "INSERT INTO picture(name, path, label) values(?,?,?)"

var sqlUpdatePrediction = "UPDATE picture SET prediction=? WHERE name=?"
var sqlUpdateLabel = "UPDATE picture SET label=? WHERE name=?"

var sqlQueryTry = "SELECT name FROM picture WHERE id=1;"

var sqlQueryName = "SELECT path, prediction, label FROM picture WHERE name=?;"


// opens a connection to the database.
func openDb() *sql.DB {
	// connect to a MySQL server as a user named localuser to a database named august2021
	db, err := sql.Open("mysql", "localuser:localuserpassword@tcp(localhost:3306)/august2021")
	CheckErr(err)
	fmt.Println("Db connection opens.")
	return db
}

// closes the connection.
func closeDb(db *sql.DB) {
	err := db.Close()
	CheckErr(err)
	fmt.Println("Db connection closes.")
}

// tests if the database functions properly.
func testDb(db *sql.DB) error {
	_, err := db.Query(sqlQueryTry)
	CheckErr(err)
	return err
}

// RecreateTable deletes the old table and rebuild a new one with default values.
func RecreateTable(db *sql.DB) {
	_, err := db.Exec(sqlDropTable)
	CheckErr(err)

	_, err = db.Exec(sqlCreateTable)
	CheckErr(err)

	_, err = db.Exec(sqlInsert, templatePNG, templatePath, nil, true)
	CheckErr(err)
}

// TryConnection tries the connection to the database and checks if it functions properly.
func TryConnection() {
	db := openDb()
	defer closeDb(db)
	err := testDb(db)
	if err != nil {
		RecreateTable(db)
		fmt.Println("Table rebuilt.")
	}else{
		fmt.Println("The database is ready.")
	}
}

// Insert inserts a new record into the database.
func Insert(name string, path string, prediction bool, label bool) {
	db := openDb()
	defer closeDb(db)

	res, err := db.Exec(sqlInsert, name, path, prediction, label)
	CheckErr(err)

	id, err := res.LastInsertId()
	CheckErr(err)

	fmt.Println("The last inserted Id is: ", id)
}

// InsertWithPrediction inserts a new record into the database.
func InsertWithPrediction(name string, path string, prediction bool) {
	db := openDb()
	defer closeDb(db)

	res, err := db.Exec(sqlInsertWithPrediction, name, path, prediction)
	CheckErr(err)

	id, err := res.LastInsertId()
	CheckErr(err)

	fmt.Println("The last inserted Id is: ", id)
}

// InsertWithLabel inserts a new record into the database.
func InsertWithLabel(name string, path string, label bool) {
	db := openDb()
	defer closeDb(db)

	res, err := db.Exec(sqlInsertWithLabel, name, path, label)
	CheckErr(err)

	id, err := res.LastInsertId()
	CheckErr(err)

	fmt.Println("The last inserted Id is: ", id)
}

// UpdatePrediction updates the prediction attribute
func UpdatePrediction(name string, prediction bool) {
	db := openDb()
	defer closeDb(db)

	res, err := db.Exec(sqlUpdatePrediction, prediction, name)
	CheckErr(err)

	n, err := res.RowsAffected()
	CheckErr(err)
	fmt.Printf("Updates Succeeds. Affected rows: %d\n", n)
}

// UpdateLabel updates the prediction attribute
func UpdateLabel(name string, label bool) {
	db := openDb()
	defer closeDb(db)

	res, err := db.Exec(sqlUpdateLabel, label, name)
	CheckErr(err)

	n, err := res.RowsAffected()
	CheckErr(err)
	fmt.Printf("Updates Succeeds. Affected rows: %d\n", n)
}

// QueryName checks if the name is in the database.
// Returns a tuple, of which the first bool means whether the name is in the database.
// The second bool is "prediction" and the third bool is "label".
// If the record with the name does not exist, the 2 booleans will always be nil.
// The last element in the tuple is a string, which is  the path attribute of
// the record, and it will be nil if the record does not exist.
func QueryName(name string) (bool, *bool, *bool, *string) {
	db := openDb()
	defer closeDb(db)
	res, err := db.Query(sqlQueryName, name)

	defer func(res *sql.Rows) {
		err := res.Close()
		CheckErr(err)
	}(res)

	CheckErr(err)

	if res.Next(){
		path := ""
		prediction := false
		label := false
		err := res.Scan(&path, &prediction, &label)
		CheckErr(err)
		return true, &prediction, &label, &path
	}else {
		return false, nil, nil, nil
	}
}
