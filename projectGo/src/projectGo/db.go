package projectGo

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)


var templatePNG = "template.PNG"

var sqlDropTable = "DROP TABLE IF EXISTS picture;"

var sqlCreateTable = "CREATE TABLE picture" + "(" +
	"id INTEGER NOT NULL AUTO_INCREMENT," +
	"name VARCHAR(100) UNIQUE NOT NULL," +
	"path VARCHAR(200) NOT NULL," +
	"b INTEGER NOT NULL," +
	"PRIMARY KEY(id)," +
	"UNIQUE (name)" + ");"

var sqlInsert = "INSERT INTO picture(name, path, b) values(?,?,?)"

var sqlQueryTry = "SELECT name FROM picture WHERE id=1;"

var sqlQueryName = "SELECT path, b FROM picture WHERE name=?;"


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

	_, err = db.Exec(sqlInsert, templatePNG, generateImgPath(templatePNG), 0)
	CheckErr(err)
}

// TryConnection tries the connection to the database and checks if it functions properly.
func TryConnection() {
	db := openDb()
	err := testDb(db)
	if err != nil {
		RecreateTable(db)
		fmt.Println("Table rebuilt.")
	}else{
		fmt.Println("The database is ready.")
	}
	closeDb(db)
}

// Insert inserts a new record into the database.
func Insert(name string, b bool) {
	db := openDb()

	res, err := db.Exec(sqlInsert, name, generateImgPath(name), mapBool2Int(b))
	CheckErr(err)

	id, err := res.LastInsertId()
	CheckErr(err)

	fmt.Println("The last inserted Id is: ", id)
}

// QueryName checks if the name is in the database.
// Returns a tuple, of which the first bool means whether the name is in the
// database and the second bool is the 'b' attribute of that record.
// If the record with the name does not exist, the second bool will always be false.
// The last element in the tuple is a string, which is  the path attribute of
// the record, and it will be "" if the record does not exist.
func QueryName(name string) (bool, bool, string) {
	db := openDb()
	defer closeDb(db)
	res, err := db.Query(sqlQueryName, name)

	defer func(res *sql.Rows) {
		err := res.Close()
		CheckErr(err)
	}(res)

	CheckErr(err)

	b := -1
	path := ""

	if res.Next(){
		err := res.Scan(&path, &b)
		CheckErr(err)
	}

	firstBool := false
	if b != -1 {
		firstBool = true
	}

	return firstBool, mapInt2Bool(b), path
}
