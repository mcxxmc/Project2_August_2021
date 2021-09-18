package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webserver/common"
)

// Db the connection pool;
// needs to be initialized by calling OpenSharedDb();
// can be closed using CloseSharedDb()
var Db *gorm.DB

// openDb opens a connection to the database.
func openDb() *gorm.DB {
	// connect to a MySQL server as a user named "localuser" to a database named august2021
	// gorm.Open() gives a connection pool to be reused frequently
	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	common.CheckErr(err)
	return db
}

// closeDb closes the connection.
func closeDb(db *gorm.DB) {
	connection, err := db.DB()
	common.CheckErr(err)
	err = connection.Close()
	common.CheckErr(err)
}

// OpenSharedDb the function to initialize the shared connection pool
func OpenSharedDb() {
	Db = openDb()
	common.Logger.Infof("Shared Db connection opens.")
}

// CloseSharedDb the function to close the shared connection pool
func CloseSharedDb() {
	closeDb(Db)
	common.Logger.Infof("Shared Db connection closes.")
}

// Tests if the database functions properly by inserting, querying and deleting a temporary record.
func testDb(db *gorm.DB) {
	hasTable := db.Migrator().HasTable(&RecordNoId{})
	if hasTable == true {
		common.Logger.Infof("db/basic.go testDb: Table found.")
	} else {
		common.Logger.Errorf("db/basic.go testDb: Error: No such table!")
	}
}

// TryConnection tries the connection to the database and checks if it functions properly.
func TryConnection() {
	db := openDb()
	defer closeDb(db)
	// Check if the table exist:
	testDb(db)
}
