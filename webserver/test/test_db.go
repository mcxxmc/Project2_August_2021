package test

import "webserver/db"

func DbConnection() {
	db.TryConnection()
}
