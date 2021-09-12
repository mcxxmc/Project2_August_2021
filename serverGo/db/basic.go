package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"serverGo/common"
)

// OpenDb opens a connection to the database.
func OpenDb() *gorm.DB {
	// connect to a MySQL server as a user named "localuser" to a database named august2021
	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	common.CheckErr(err)
	// fmt.Println("Db connection opens.")
	return db
}

// CloseDb closes the connection.
func CloseDb(db *gorm.DB) {
	connection, err := db.DB()
	common.CheckErr(err)
	err = connection.Close()
	common.CheckErr(err)
	// fmt.Println("Db connection closes.")
}

// Tests if the database functions properly by inserting, querying and deleting a temporary record.
func testDb(db *gorm.DB) {
	hasTable := db.Migrator().HasTable(&RecordNoId{})
	if hasTable == true {
		fmt.Println("db/basic.go testDb: Table found.")
	} else {
		fmt.Println("db/basic.go testDb: Error: No such table!")
	}
}

// TryConnection tries the connection to the database and checks if it functions properly.
func TryConnection() {
	db := OpenDb()
	defer CloseDb(db)
	// Check if the table exist:
	testDb(db)
}
