package main

import (
	"go.uber.org/zap"
	"webserver/db"
	"webserver/tf_implement"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()  // flushes buffer, if any
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	db.OpenSharedDb()
	defer db.CloseSharedDb()
	tf_implement.StartServer()
}
