package main

import (
	"webserver/common"
	"webserver/db"
	"webserver/tf_implement"
)

func main() {
	common.InitLog()
	db.OpenSharedDb()
	defer db.CloseSharedDb()
	tf_implement.StartServer()
}
